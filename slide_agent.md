# Slide Agent Interface Design
*Building AI Agents for PowerPoint Editing with UNO API*

## Overview

Following the pattern from Thorsten Ball's agent building approach, we need to create a set of **tools** that allow an AI agent to intelligently edit PowerPoint presentations using LibreOffice's UNO API.

Just like the code-editing agent had `read_file`, `list_files`, and `edit_file`, our slide agent needs equivalent tools for presentations.

## Core Tool Set

### 1. `list_slides` - Presentation Overview
**Purpose**: Get an overview of all slides in a presentation
**Usage**: "What slides are in this presentation?" or "Show me the structure"

```json
{
  "name": "list_slides",
  "description": "List all slides in a presentation with their basic information including slide number, title, and layout type.",
  "input_schema": {
    "presentation_path": "string - Path to the .pptx file"
  },
  "returns": "JSON array of slide objects with id, title, layout, slide_number"
}
```

**Example Output**:
```json
[
  {"slide_number": 1, "title": "Hello", "layout": "Title Slide", "text_shapes": 2},
  {"slide_number": 2, "title": "Welcome", "layout": "Content Slide", "text_shapes": 3}
]
```

### 2. `read_slide` - Slide Content Inspection
**Purpose**: Get detailed content from a specific slide with shape-level detail
**Usage**: "What's on slide 3?" or "Read the content of the first slide"

```json
{
  "name": "read_slide", 
  "description": "Read all text content from a specific slide with detailed shape information including shape indices, types, and content.",
  "input_schema": {
    "presentation_path": "string - Path to the .pptx file",
    "slide_number": "integer - Slide number (1-based)"
  },
  "returns": "Structured content with shape indices for precise targeting"
}
```

**Example Output**:
```json
{
  "slide_number": 1,
  "shapes": [
    {
      "shape_index": 0,
      "shape_type": "title", 
      "text": "Hello",
      "description": "Main slide title"
    },
    {
      "shape_index": 1,
      "shape_type": "content",
      "text": "• First bullet point\n• Second bullet point\n• Third bullet point",
      "description": "Content text box with bullet points",
      "bullet_points": [
        {"index": 0, "text": "First bullet point"},
        {"index": 1, "text": "Second bullet point"}, 
        {"index": 2, "text": "Third bullet point"}
      ]
    },
    {
      "shape_index": 2,
      "shape_type": "text_box",
      "text": "Additional notes here",
      "description": "Additional text box"
    }
  ]
}
```

### 3. `edit_slide_text` - Content Modification
**Purpose**: Modify text content on a slide with precise targeting
**Usage**: "Change the title to 'Goodbye'" or "Update the third bullet point"

```json
{
  "name": "edit_slide_text",
  "description": "Edit text content on a slide by targeting specific shapes or elements. Can target by shape index, shape type, or bullet point index for precise editing.",
  "input_schema": {
    "presentation_path": "string - Path to the .pptx file", 
    "slide_number": "integer - Slide number (1-based)",
    "target_type": "string - How to target: 'shape_index', 'shape_type', 'bullet_point', or 'text_replace'",
    "target_value": "string/integer - Shape index (0,1,2...), shape type ('title','content'), bullet index (0,1,2...), or text to find",
    "new_text": "string - New text content",
    "old_text": "string - (Optional) For text_replace mode, the exact text to replace"
  },
  "returns": "Success message or error"
}
```

**Usage Examples**:
```json
// Change slide title
{"target_type": "shape_type", "target_value": "title", "new_text": "New Title"}

// Change specific shape by index
{"target_type": "shape_index", "target_value": 1, "new_text": "Updated content"}

// Change third bullet point (index 2)
{"target_type": "bullet_point", "target_value": 2, "new_text": "New third bullet"}

// Traditional text replacement
{"target_type": "text_replace", "old_text": "Hello", "new_text": "Goodbye"}
```

### 4. `add_slide` - Slide Creation
**Purpose**: Add new slides to the presentation
**Usage**: "Add a new slide with title 'Conclusion'" or "Create a blank slide"

```json
{
  "name": "add_slide",
  "description": "Add a new slide to the presentation with optional initial content.",
  "input_schema": {
    "presentation_path": "string - Path to the .pptx file",
    "position": "integer - Position to insert (optional, defaults to end)",
    "layout": "string - Slide layout type (optional, defaults to 'blank')",
    "title": "string - Initial title text (optional)"
  },
  "returns": "New slide number and confirmation"
}
```

### 5. `delete_slide` - Slide Removal
**Purpose**: Remove slides from the presentation
**Usage**: "Delete slide 3" or "Remove the last slide"

```json
{
  "name": "delete_slide",
  "description": "Delete a slide from the presentation.",
  "input_schema": {
    "presentation_path": "string - Path to the .pptx file",
    "slide_number": "integer - Slide number to delete (1-based)"
  },
  "returns": "Confirmation message"
}
```

### 6. `export_slides` - Visual Output
**Purpose**: Generate visual representation of slides (for verification/preview)
**Usage**: "Show me what the slides look like" or "Export slide 2 as an image"

```json
{
  "name": "export_slides",
  "description": "Export slides as JPEG images for preview or verification.",
  "input_schema": {
    "presentation_path": "string - Path to the .pptx file",
    "slide_numbers": "array of integers - Specific slides to export (optional, defaults to all)",
    "output_dir": "string - Directory to save images (optional, defaults to 'slides/')"
  },
  "returns": "Array of exported image file paths"
}
```

## Example Agent Interaction Flow

```
User: "I have a presentation called demo.pptx. Can you change the title of the first slide from 'Hello' to 'Welcome Everyone' and update the third bullet point to say 'New conclusion' and then show me what it looks like?"

Agent thinks: I need to:
1. First understand what's in the presentation
2. Modify the title and specific bullet point
3. Export to see the result

tool: list_slides({"presentation_path": "demo.pptx"})
→ Shows 2 slides, first slide has title "Hello"

tool: read_slide({"presentation_path": "demo.pptx", "slide_number": 1})
→ Shows slide structure with shape indices and bullet points

tool: edit_slide_text({
  "presentation_path": "demo.pptx", 
  "slide_number": 1,
  "target_type": "shape_type",
  "target_value": "title", 
  "new_text": "Welcome Everyone"
})
→ "Successfully updated slide 1 title"

tool: edit_slide_text({
  "presentation_path": "demo.pptx", 
  "slide_number": 1,
  "target_type": "bullet_point",
  "target_value": 2,
  "new_text": "New conclusion"
})
→ "Successfully updated bullet point 3 on slide 1"

tool: export_slides({
  "presentation_path": "demo.pptx",
  "slide_numbers": [1]
})
→ ["slides/slide-000.jpg"]

Agent: "I've successfully changed the title of the first slide from 'Hello' to 'Welcome Everyone' and updated the third bullet point to 'New conclusion'. The updated slide has been exported as an image for you to review."
```

## Implementation Architecture

```
┌─────────────────────────────┐
│ AI Agent (Go + Anthropic)   │
│  • Tool definitions         │
│  • Conversation loop        │
│  • Tool execution           │
└─────────────┬───────────────┘
              │ Calls tools
              ▼
┌─────────────────────────────┐
│ UNO Tool Functions (Go)     │
│  • list_slides()            │
│  • read_slide()             │
│  • edit_slide_text()        │
│  • add_slide()              │
│  • export_slides()          │
└─────────────┬───────────────┘
              │ Spawns Python
              ▼
┌─────────────────────────────┐
│ Python UNO Scripts          │
│  • uno_list_slides.py       │
│  • uno_read_slide.py        │
│  • uno_edit_slide.py        │
│  • uno_add_slide.py         │
└─────────────┬───────────────┘
              │ UNO API calls
              ▼
┌─────────────────────────────┐
│ LibreOffice Headless        │
│  • Port 8100 UNO socket     │
│  • Document manipulation    │
│  • Persistent service       │
└─────────────────────────────┘
```

## Tool Design Principles

Following the article's approach:

1. **Simple & Focused**: Each tool does one thing well
2. **String-Based**: AI models work best with text-based operations
3. **Forgiving**: Good error messages and graceful failure
4. **Composable**: Tools can be chained naturally by the AI
5. **Observable**: Clear feedback about what happened

## Scaling Challenges & Strategy

**The Complexity Problem**: As we add features (colors, fonts, animations, transitions, layouts, images, charts), we risk tool explosion:

```
❌ BAD - Tool Explosion:
edit_slide_text, edit_text_color, edit_text_font, edit_text_size,
edit_slide_background, edit_slide_transition, add_animation,
edit_animation_timing, add_image, resize_image, move_image,
add_chart, edit_chart_data, edit_chart_colors...
```

**Solution: Hierarchical Tool Design**

### Core Tools (MVP - What We're Building)
- `list_slides` - Structure
- `read_slide` - Content inspection  
- `edit_slide_text` - Text content
- `add_slide` / `delete_slide` - Structure modification
- `export_slides` - Verification

### Future Tool Groups (When Needed)

**Visual Tools** (Phase 2):
```json
{
  "name": "edit_slide_style",
  "description": "Modify visual properties like colors, fonts, backgrounds",
  "input_schema": {
    "style_type": "background|text_color|font|font_size",
    "target_type": "slide|shape_index|shape_type|bullet_point",
    "style_value": "red|Arial|18px|etc"
  }
}
```

**Layout Tools** (Phase 3):
```json
{
  "name": "modify_slide_layout", 
  "description": "Change slide layouts, move elements, resize shapes",
  "input_schema": {
    "layout_action": "change_template|move_shape|resize_shape",
    "target": "shape_index|coordinates",
    "new_value": "Title Slide|{x:100,y:200}|{width:400,height:300}"
  }
}
```

**Media Tools** (Phase 4):
```json
{
  "name": "manage_slide_media",
  "description": "Add, edit, or remove images, charts, and media",
  "input_schema": {
    "media_action": "add_image|add_chart|edit_chart_data|remove_media",
    "media_params": "varies by action type"
  }
}
```

### Design Strategy for Scale

1. **Start Simple**: Build the 6 core tools first, get them working perfectly
2. **Semantic Grouping**: Group related functionality (text, visual, layout, media)
3. **Consistent Targeting**: Use same targeting system (`shape_index`, `shape_type`, etc.) across all tools
4. **Progressive Enhancement**: Add complexity only when core use cases are solid
5. **AI-Friendly Abstractions**: Prefer "make this slide professional" over "set font to Arial 18px, color to #333333, background to gradient..."

### The 80/20 Rule

**80% of slide editing is**:
- Changing text content ✅ (we have this)
- Adding/removing slides ✅ (we have this) 
- Basic visual formatting (colors, fonts)
- Adding simple images

**20% is advanced**:
- Complex animations
- Custom layouts
- Advanced charts
- Video/audio

**Strategy**: Perfect the 80% first, then add the 20% thoughtfully.

## Why This Works

Just like in the code-editing example:

- **No complex prompting needed**: AI figures out when to use tools
- **Natural chaining**: AI will `list_slides` → `read_slide` → `edit_slide_text` → `export_slides`
- **Minimal code**: Each tool is ~20-50 lines of Go + Python
- **Reliable**: UNO API is stable and well-documented

## Next Steps

1. Implement each tool as a Go function that calls Python UNO scripts
2. Create the Agent framework (copy from article, adapt for slides)
3. Add tool definitions with proper JSON schemas
4. Test with simple commands like "change title of slide 1"
5. Iterate based on what the AI actually tries to do

The magic isn't in complex prompting or clever algorithms - it's in giving the AI the right set of simple, composable tools and letting it figure out how to use them.
