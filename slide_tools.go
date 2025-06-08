package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// ListSlidesDefinition defines the list_slides tool
var ListSlidesDefinition = ToolDefinition{
	Name: "list_slides",
	Description: `List all slides in a PowerPoint presentation with basic information.
	
Use this tool to get an overview of the presentation structure, including slide numbers, titles, and layout information. This is typically the first tool to use when working with a presentation.`,
	InputSchema: ListSlidesInputSchema,
	Function:    ListSlides,
}

type ListSlidesInput struct {
	PresentationPath string `json:"presentation_path" jsonschema_description:"Path to the PowerPoint (.pptx) file"`
}

var ListSlidesInputSchema = GenerateSchema[ListSlidesInput]()

func ListSlides(input json.RawMessage) (string, error) {
	listSlidesInput := ListSlidesInput{}
	err := json.Unmarshal(input, &listSlidesInput)
	if err != nil {
		return "", fmt.Errorf("failed to parse input: %v", err)
	}

	if listSlidesInput.PresentationPath == "" {
		return "", fmt.Errorf("presentation_path is required")
	}

	fmt.Printf("Listing slides in: %s\n", listSlidesInput.PresentationPath)

	// Call Python UNO script
	cmd := exec.Command("python3", "uno_list_slides.py", listSlidesInput.PresentationPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to list slides: %v\nOutput: %s", err, string(output))
	}

	// Validate that the output is valid JSON
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("invalid JSON output from UNO script: %v", err)
	}

	return string(output), nil
}

// ReadSlideDefinition defines the read_slide tool
var ReadSlideDefinition = ToolDefinition{
	Name: "read_slide",
	Description: `Read detailed content from a specific slide including all text shapes and their content.

Use this tool to get detailed information about a specific slide's content, including shape indices, types, and text content. This is essential for understanding slide structure before making edits.`,
	InputSchema: ReadSlideInputSchema,
	Function:    ReadSlide,
}

type ReadSlideInput struct {
	PresentationPath string `json:"presentation_path" jsonschema_description:"Path to the PowerPoint (.pptx) file"`
	SlideNumber      int    `json:"slide_number" jsonschema_description:"Slide number to read (1-based indexing)"`
}

var ReadSlideInputSchema = GenerateSchema[ReadSlideInput]()

func ReadSlide(input json.RawMessage) (string, error) {
	readSlideInput := ReadSlideInput{}
	err := json.Unmarshal(input, &readSlideInput)
	if err != nil {
		return "", fmt.Errorf("failed to parse input: %v", err)
	}

	if readSlideInput.PresentationPath == "" {
		return "", fmt.Errorf("presentation_path is required")
	}

	if readSlideInput.SlideNumber < 1 {
		return "", fmt.Errorf("slide_number must be 1 or greater")
	}

	fmt.Printf("Reading slide %d from: %s\n", readSlideInput.SlideNumber, readSlideInput.PresentationPath)

	// Call Python UNO script
	cmd := exec.Command("python3", "uno_read_slide.py", readSlideInput.PresentationPath, fmt.Sprintf("%d", readSlideInput.SlideNumber))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to read slide: %v\nOutput: %s", err, string(output))
	}

	// Validate that the output is valid JSON
	var result interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("invalid JSON output from UNO script: %v", err)
	}

	return string(output), nil
}

// EditSlideTextDefinition defines the edit_slide_text tool
var EditSlideTextDefinition = ToolDefinition{
Name: "edit_slide_text",
Description: `Edit text content on a slide by targeting specific shapes or elements. 

Can target by shape index, shape type, or replace specific text. This tool allows precise editing of slide content including titles, text boxes, and bullet points.

Target types:
- "shape_index": Edit specific shape by index (0, 1, 2, ...)
- "shape_type": Edit by type ("title", "content", "text_box")  
- "text_replace": Replace specific text (requires old_text)
- "bullet_point": Edit specific bullet point by index`,
InputSchema: EditSlideTextInputSchema,
Function:    EditSlideText,
}

type EditSlideTextInput struct {
PresentationPath string `json:"presentation_path" jsonschema_description:"Path to the PowerPoint (.pptx) file"`
SlideNumber      int    `json:"slide_number" jsonschema_description:"Slide number to edit (1-based indexing)"`
TargetType       string `json:"target_type" jsonschema_description:"How to target: 'shape_index', 'shape_type', 'bullet_point', or 'text_replace'"`
TargetValue      string `json:"target_value" jsonschema_description:"Shape index (0,1,2...), shape type ('title','content','text_box'), bullet index, or text to find"`
NewText          string `json:"new_text" jsonschema_description:"New text content to set"`
OldText          string `json:"old_text,omitempty" jsonschema_description:"(Optional) For text_replace mode, the exact text to replace"`
}

var EditSlideTextInputSchema = GenerateSchema[EditSlideTextInput]()

func EditSlideText(input json.RawMessage) (string, error) {
editInput := EditSlideTextInput{}
err := json.Unmarshal(input, &editInput)
if err != nil {
return "", fmt.Errorf("failed to parse input: %v", err)
}

if editInput.PresentationPath == "" {
return "", fmt.Errorf("presentation_path is required")
}

if editInput.SlideNumber < 1 {
return "", fmt.Errorf("slide_number must be 1 or greater")
}

if editInput.TargetType == "" {
return "", fmt.Errorf("target_type is required")
}

if editInput.TargetValue == "" {
return "", fmt.Errorf("target_value is required")
}

if editInput.NewText == "" {
return "", fmt.Errorf("new_text is required")
}

if editInput.TargetType == "text_replace" && editInput.OldText == "" {
return "", fmt.Errorf("old_text is required for text_replace mode")
}

fmt.Printf("Editing slide %d: %s=%s -> '%s'\n", 
editInput.SlideNumber, editInput.TargetType, editInput.TargetValue, editInput.NewText)

// Build command arguments
args := []string{
"uno_edit_slide.py",
editInput.PresentationPath,
fmt.Sprintf("%d", editInput.SlideNumber),
editInput.TargetType,
editInput.TargetValue,
editInput.NewText,
}

// Add old_text if provided
if editInput.OldText != "" {
args = append(args, editInput.OldText)
}

// Call Python UNO script
cmd := exec.Command("python3", args...)
output, err := cmd.CombinedOutput()
if err != nil {
return "", fmt.Errorf("failed to edit slide: %v\nOutput: %s", err, string(output))
}

// Validate that the output is valid JSON
var result interface{}
if err := json.Unmarshal(output, &result); err != nil {
return "", fmt.Errorf("invalid JSON output from UNO script: %v", err)
}

return string(output), nil
}

// ExportSlidesDefinition defines the export_slides tool
var ExportSlidesDefinition = ToolDefinition{
Name: "export_slides",
Description: `Export slides as JPEG images for preview or verification.

Use this tool to generate visual representations of slides, especially useful after making edits to verify changes. Can export all slides or specific slides.`,
InputSchema: ExportSlidesInputSchema,
Function:    ExportSlides,
}

type ExportSlidesInput struct {
PresentationPath string `json:"presentation_path" jsonschema_description:"Path to the PowerPoint (.pptx) file"`
SlideNumbers     []int  `json:"slide_numbers,omitempty" jsonschema_description:"Specific slides to export (optional, defaults to all slides)"`
OutputDir        string `json:"output_dir,omitempty" jsonschema_description:"Directory to save images (optional, defaults to 'slides/')"`
}

var ExportSlidesInputSchema = GenerateSchema[ExportSlidesInput]()

func ExportSlides(input json.RawMessage) (string, error) {
exportInput := ExportSlidesInput{}
err := json.Unmarshal(input, &exportInput)
if err != nil {
return "", fmt.Errorf("failed to parse input: %v", err)
}

if exportInput.PresentationPath == "" {
return "", fmt.Errorf("presentation_path is required")
}

// Set default output directory
outputDir := exportInput.OutputDir
if outputDir == "" {
outputDir = "slides"
}

fmt.Printf("Exporting slides from: %s to %s/\n", exportInput.PresentationPath, outputDir)

// Use our existing conversion function
slides, err := ConvertPPTXToJPEG(exportInput.PresentationPath, outputDir)
if err != nil {
return "", fmt.Errorf("failed to export slides: %v", err)
}

// Filter slides if specific slide numbers were requested
var filteredSlides []string
if len(exportInput.SlideNumbers) > 0 {
slideMap := make(map[int]bool)
for _, num := range exportInput.SlideNumbers {
slideMap[num-1] = true // Convert to 0-based indexing
}

for i, slide := range slides {
if slideMap[i] {
filteredSlides = append(filteredSlides, slide)
}
}
slides = filteredSlides
}

result := map[string]interface{}{
"success":     true,
"slide_count": len(slides),
"slides":      slides,
"output_dir":  outputDir,
}

resultJSON, _ := json.Marshal(result)
return string(resultJSON), nil
}


