# 🗂️ Proof‑of‑Concept Plan — “Wails + LibreOffice Headless”

| Layer                      | What It Does                                                                                          | Main Tech                     | Key Files                                     |
| -------------------------- | ----------------------------------------------------------------------------------------------------- | ----------------------------- | --------------------------------------------- |
| **UI (Renderer)**          | Displays slides, user controls, basic “edit” buttons                                                  | React (via Wails), HTML / CSS | `frontend/src/App.tsx`, `SlideViewer.tsx`     |
| **App Shell**              | Starts UI, exposes Go methods to JS, emits events                                                     | **Wails** (v3)                | `wails.json`, `main.go`                       |
| **Go Service**             | 1) Reads/receives a `.pptx` 2) Calls conversion / edit helpers 3) Serves images or Base‑64 back to UI | Go 1.22                       | `backend/convert.go`, `backend/uno_bridge.go` |
| **Python UNO Helper**      | Connects to the running headless LibreOffice, applies edits, triggers export per slide                | Python 3, PyUNO               | `uno/edit_slide.py`, `uno/export_slide.py`    |
| **LibreOffice (Headless)** | Renders & saves PDF/images; stays alive on port 8100 for UNO                                          | LibreOffice 7.x               | Installed binary `soffice`                    |

---

## 1 · Environment Setup

```bash
# 1. Install LibreOffice + headless extras
sudo apt-get install libreoffice libreoffice-impress libreoffice-headless

# 2. Install Python + PyUNO bindings (Ubuntu example)
sudo apt-get install python3-uno

# 3. Get ImageMagick (PDF→PNG)
sudo apt-get install imagemagick

# 4. Wails + Go
go install github.com/wailsapp/wails/v3/cmd/wails@latest
```

---

## 2 · Runtime Architecture

```text
┌────────────────────────────┐
│ React (Wails Renderer)     │
│  • SlideViewer component   │
│  • “Edit Title” button     │
└──────┬─────────▲───────────┘
       ▼         │  JSON RPC
┌────────────────────────────┐
│ Go Backend (Wails)         │
│  convert.go                │
│  uno_bridge.go             │
│  • spawn libreoffice (pdf) │
│  • spawn convert (png)     │
│  • call python UNO script  │
└──────┬─────────▲───────────┘
       ▼         │  stdin/stdout
┌────────────────────────────┐
│ Python (pyuno)             │
│  • edit_slide.py           │
│  • export_slide.py         │
└──────┬─────────▲───────────┘
       ▼ socket  │ UNO API
┌────────────────────────────┐
│ LibreOffice headless       │
│  --accept='socket,port=8100│
└────────────────────────────┘
```

---

## 3 · Execution Flow

| #     | Action                       | Behind the Scenes                                                                                                                                              |
| ----- | ---------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **0** | App starts                   | Wails launches; Go spawns LibreOffice headless once: `soffice --headless --invisible --accept="socket,host=127.0.0.1,port=8100;urp;StarOffice.ServiceManager"` |
| **1** | User opens a `.pptx`         | Go runs `soffice --headless --convert-to pdf <file>` → gets `slides.pdf`; then `convert slides.pdf slide-%d.png`                                               |
| **2** | UI displays slides           | Go returns array of `slide‑0.png …`; React renders them in `<img>` tags or a canvas                                                                            |
| **3** | User clicks **Edit Slide 2** | Frontend calls `backend.EditSlide(2,"New title")`                                                                                                              |
| **4** | Go → Python                  | Go spawns `python3 uno/edit_slide.py 2 "New title" file:///<pptx>`                                                                                             |
| **5** | Python UNO does edit         | Script connects on port 8100, loads doc, picks page 2, edits text box, saves doc                                                                               |
| **6** | Go re‑exports slide 2 only   | Python triggers UNO “export graphic” filter or Go reruns PDF→PNG for that page                                                                                 |
| **7** | Frontend refreshes image     | Go fires Wails event `slideUpdated` with new Base‑64; React swaps only that slide                                                                              |

---

## 4 · Key Code Snippets

**convert.go**

```go
func ConvertToPNG(pptPath string) ([]string, error) {
    tmpDir := os.MkdirTemp("", "slides")
    cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf",
                        "--outdir", tmpDir, pptPath)
    if err := cmd.Run(); err != nil { return nil, err }
    pdf := filepath.Join(tmpDir, strings.TrimSuffix(filepath.Base(pptPath), ".pptx")+".pdf")
    exec.Command("convert", "-density", "150", pdf, filepath.Join(tmpDir, "slide-%d.png")).Run()
    files, _ := filepath.Glob(filepath.Join(tmpDir, "slide-*.png"))
    return files, nil
}
```

**edit_slide.py (PyUNO skeleton)**

```python
import uno, sys
slide_idx, new_text, doc_url = int(sys.argv[1]), sys.argv[2], sys.argv[3]

ctx = uno.getComponentContext()
sm  = ctx.ServiceManager
resolver = sm.createInstanceWithContext("com.sun.star.bridge.UnoUrlResolver", ctx)
lo = resolver.resolve("uno:socket,host=localhost,port=8100;urp;StarOffice.ComponentContext")
desktop = lo.ServiceManager.createInstanceWithContext("com.sun.star.frame.Desktop", lo)
doc = desktop.loadComponentFromURL(doc_url, "_blank", 0, ())
page = doc.getDrawPages().getByIndex(slide_idx)
shape = page.getByIndex(0)                 # assume first shape is title
shape.setString(new_text)
doc.store()                                # save .pptx
```

---

## 5 · Directory Skeleton

```
/poc
 ├─ frontend/              (React via Wails)
 ├─ backend/
 │   ├─ convert.go
 │   ├─ uno_bridge.go
 ├─ uno/
 │   ├─ edit_slide.py
 │   └─ export_slide.py
 └─ assets/
     └─ sample.pptx
```

---

## 6 · Build & Run

```bash
# Run the LibreOffice service in the background
soffice --headless --invisible --norestore \
  --accept="socket,host=127.0.0.1,port=8100;urp;StarOffice.ServiceManager" &

# Start the Wails app (dev mode)
wails dev

# → Open the desktop window, load sample.pptx, edit, watch slide refresh
```

---

## 7 · Stretch‑Goals After POC

| Goal                   | Notes                                                                    |
| ---------------------- | ------------------------------------------------------------------------ |
| **Persistent preview** | Keep PNGs cached, only regenerate edited slide                           |
| **Undo/Redo**          | Version the `.pptx` after each edit                                      |
| **AI agent**           | Expose gRPC/REST in Go; pass prompts to agent → call `edit_slide.py`     |
| **Packaging**          | `wails build -upx` then bundle LibreOffice, Python, scripts in installer |

---

## ✅ What This POC Proves

1. **Render fidelity** – LibreOffice → PDF → PNG pipeline is good enough for slide previews.
2. **Round‑trip editing** – Text changes via UNO are reflected live in the UI.
3. **Wails interaction** – Front‑end↔Go IPC latency is acceptable for “AI‑powered editor”.

Once this all works locally, you can iterate on performance (persistent LO process, caching) and eventually decide whether to stay on Wails or migrate to a Python/Electron stack for deeper UNO integration.
