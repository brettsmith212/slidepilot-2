package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <path-to-pptx-file>")
		os.Exit(1)
	}

	pptxPath := os.Args[1]
	fmt.Printf("Converting PPTX: %s\n", pptxPath)

	slides, err := ConvertPPTXToJPEG(pptxPath)
	if err != nil {
		log.Fatalf("Conversion failed: %v", err)
	}

	fmt.Printf("Successfully converted to %d slides:\n", len(slides))
	for i, slide := range slides {
		fmt.Printf("  Slide %d: %s\n", i+1, slide)
	}
}
