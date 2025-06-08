# SlidePilot - AI Slide Editing Agent

## Overview

Successfully built a **complete AI slide editing agent** following Thorsten Ball's agent building approach from "How to Build an Agent". The agent can intelligently edit PowerPoint presentations using LibreOffice's UNO API with visual verification.

## Core Architecture

**4 Essential Tools** (following the article's pattern):
- `list_slides` - Get presentation overview
- `read_slide` - Detailed slide content analysis  
- `edit_slide_text` - Modify slide content with precise targeting
- `export_slides` - Generate JPEG images for visual verification

**Technology Stack**:
- **Go + Anthropic SDK** - Agent framework and conversation loop
- **Python + UNO API** - LibreOffice automation for slide editing
- **LibreOffice Headless** - Document processing engine
- **ImageMagick** - PDF to JPEG conversion for visual output

## Key Files

### Core Agent Files
- `main.go` - AI agent with conversation loop and tool execution
- `slide_tools.go` - Tool definitions and Go implementations
- `libreoffice_service.go` - LibreOffice headless service management

### UNO API Scripts
- `uno_list_slides.py` - Extract slide overview information
- `uno_read_slide.py` - Get detailed slide content with shape indices
- `uno_edit_slide.py` - Edit slide text using precise targeting

### Legacy Files
- `simple_editor.go` - Original POC (renamed from main.go)
- `converter.go` - JPEG conversion utilities

## How to Use the Agent

### 1. Setup Requirements
```bash
# Required packages
sudo apt-get install libreoffice libreoffice-headless python3-uno imagemagick

# Set API key
export ANTHROPIC_API_KEY="your-key-here"
```

### 2. Start the Agent
```bash
# Build and run
go build
./slide-agent

# Agent will start LibreOffice headless automatically
# Chat interface: "Chat with Claude - Slide Agent (use 'ctrl-c' to quit)"
```

### 3. Example Commands
```
"What slides are in original_ppt.pptx?"
"Change the title of slide 1 to 'New Title'"  
"Show me what slide 1 looks like after editing"
"Read the detailed content of slide 2"
```

## Testing Instructions

### 🚨 IMPORTANT: Always Use Copies for Testing
```bash
# ALWAYS copy original before testing
cp original_ppt.pptx test_example.pptx

# Then use the copy in commands
"Edit test_example.pptx and change the title to 'Test Title'"
```

### Sample Test Workflow
```bash
# 1. Create test copy
cp original_ppt.pptx test_workflow.pptx

# 2. Start agent
./slide-agent

# 3. Test complete workflow
"Change slide 1 title in test_workflow.pptx from 'Hello' to 'AI Test' and show me the result"

# 4. Verify outputs
ls slides/  # Check for slide-000.jpg, slide-001.jpg
```

### Verification Commands
```bash
# Check slide content directly
python3 uno_read_slide.py test_file.pptx 1

# View generated images
ls -la slides/
```

## Tool Capabilities

### `list_slides`
- **Purpose**: Get presentation overview
- **Output**: Slide count, titles, text shape counts
- **Usage**: "What slides are in this presentation?"

### `read_slide` 
- **Purpose**: Detailed content analysis with shape indices
- **Output**: Shape-by-shape breakdown with targeting info
- **Usage**: "What's the detailed content of slide 2?"

### `edit_slide_text`
- **Purpose**: Precise text editing using multiple targeting modes
- **Modes**: 
  - `shape_index` - Target specific shape (0, 1, 2...)
  - `shape_type` - Target by type ("title", "content", "text_box")
  - `text_replace` - Replace specific text
  - `bullet_point` - Edit individual bullet points
- **Usage**: "Change the title to 'New Title'"

### `export_slides`
- **Purpose**: Visual verification via JPEG export
- **Output**: slide-000.jpg, slide-001.jpg in slides/ directory
- **Usage**: "Show me what the slides look like"

## Key Learnings

### 1. Agent Building Principles (from Article)
- **Simple tools work best** - 4 focused tools vs. 20+ complex ones
- **AI chains tools naturally** - No complex prompting needed
- **String-based operations** - Models excel at text manipulation
- **Tool composition** - list → read → edit → export workflow

### 2. UNO API Insights
- **Shape indexing** - Each text element gets precise index for targeting
- **Persistent service** - Keep LibreOffice running on port 8100 for performance
- **Hidden loading** - Use `PropertyValue("Hidden", 0, True, 0)` for headless operation
- **Text replacement** - Most reliable editing method for AI agents

### 3. Visual Verification Critical
- **JPEG export** - Essential for confirming edits worked
- **ImageMagick pipeline** - PPTX → PDF → JPEG for compatibility
- **File size changes** - Different content = different file sizes (verification indicator)

## Error Handling & Debugging

### Common Issues
1. **LibreOffice not running** - Agent auto-starts service on port 8100
2. **File not found** - Always use full/relative paths consistently  
3. **Edit errors** - Check with `uno_read_slide.py` for verification
4. **Permission issues** - Ensure files aren't read-only

### Debug Commands
```bash
# Check LibreOffice service
ss -ln | grep 8100

# Test UNO scripts directly
python3 uno_list_slides.py test_file.pptx
python3 uno_read_slide.py test_file.pptx 1
python3 uno_edit_slide.py test_file.pptx 1 shape_type title "Debug Title"

# Build and test
go build
```

## Build Commands

```bash
# Development build
go build

# Clean rebuild
go clean && go build

# Run simple editor (legacy)
# Note: simple_editor.go has runSimpleEditor() function, not main()
```

## Future Enhancements

Following the article's 80/20 rule, next tools to add:
- `add_slide` - Create new slides
- `delete_slide` - Remove slides  
- `edit_slide_style` - Colors, fonts, formatting
- `manage_slide_media` - Images, charts, media

## Success Metrics

✅ **Agent works like code-editing example** - Natural tool chaining  
✅ **Visual verification** - JPEG export confirms changes  
✅ **Precise targeting** - Shape-level editing accuracy  
✅ **Error recovery** - Continues despite minor issues  
✅ **Simple usage** - Natural language commands work  

## Repository Structure

```
slidepilot-2/
├── main.go              # AI agent (primary program)
├── slide_tools.go       # Tool definitions  
├── uno_*.py            # Python UNO scripts
├── libreoffice_service.go # Service management
├── original_ppt.pptx   # Test file (preserve!)
├── slides/             # JPEG output directory
├── simple_editor.go    # Legacy POC
└── AGENT.md           # This file
```

This represents a **complete, working AI slide editing agent** built in <400 lines of code, proving the article's core thesis: "It's an LLM, a loop, and enough tokens. The rest is elbow grease."
