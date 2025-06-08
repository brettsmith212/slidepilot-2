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
