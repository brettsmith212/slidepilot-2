package main

import (
	"fmt"
	"log"
	"os"
)

func runSimpleEditor() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <path-to-pptx-file>")
		os.Exit(1)
	}

	pptxPath := os.Args[1]
	fmt.Printf("Processing PPTX: %s\n", pptxPath)

	// Step 1: Ensure LibreOffice headless service is running
	if err := StartLibreOfficeHeadless(); err != nil {
		log.Fatalf("Failed to start LibreOffice service: %v", err)
	}

	// Step 2: Edit the first slide title from "Hello" to "Goodbye"
	if err := EditSlideTitle(pptxPath, 0, "Goodbye"); err != nil {
		log.Fatalf("Edit failed: %v", err)
	}

	// Step 3: Convert to JPEG slides
	slides, err := ConvertPPTXToJPEG(pptxPath)
	if err != nil {
		log.Fatalf("Conversion failed: %v", err)
	}

	fmt.Printf("Successfully converted to %d slides:\n", len(slides))
	for i, slide := range slides {
		fmt.Printf("  Slide %d: %s\n", i+1, slide)
	}
}
