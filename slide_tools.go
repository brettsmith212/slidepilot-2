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
