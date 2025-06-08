package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EditSlideTitle edits the title of a specific slide using Python UNO script
func EditSlideTitle(pptxPath string, slideIndex int, newTitle string) error {
	fmt.Printf("Editing slide %d title to '%s'...\n", slideIndex, newTitle)
	
	cmd := exec.Command("python3", "edit_slide.py", pptxPath, 
		fmt.Sprintf("%d", slideIndex), newTitle)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("edit failed: %v\nOutput: %s", err, string(output))
	}
	
	fmt.Printf("Edit output: %s", string(output))
	return nil
}

// ConvertPPTXToJPEG converts a PPTX file to JPEG slides using LibreOffice and ImageMagick
func ConvertPPTXToJPEG(pptxPath string) ([]string, error) {
	// Create slides output directory
	slidesDir := "slides"
	if err := os.MkdirAll(slidesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create slides directory: %v", err)
	}

	// Create temporary directory for PDF
	tmpDir, err := os.MkdirTemp("", "slidepilot-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Step 1: Convert PPTX to PDF using LibreOffice headless
	fmt.Println("Converting PPTX to PDF...")
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf", 
		"--outdir", tmpDir, pptxPath)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("LibreOffice conversion failed: %v", err)
	}

	// Find the generated PDF file
	baseName := strings.TrimSuffix(filepath.Base(pptxPath), ".pptx")
	pdfPath := filepath.Join(tmpDir, baseName+".pdf")
	
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("PDF file not found at %s", pdfPath)
	}

	// Step 2: Convert PDF to JPEG using ImageMagick
	fmt.Println("Converting PDF to JPEG slides...")
	outputPattern := filepath.Join(slidesDir, "slide-%03d.jpg")
	cmd = exec.Command("convert", "-density", "150", pdfPath, outputPattern)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ImageMagick conversion failed: %v", err)
	}

	// Find all generated JPEG files
	jpegFiles, err := filepath.Glob(filepath.Join(slidesDir, "slide-*.jpg"))
	if err != nil {
		return nil, fmt.Errorf("failed to find JPEG files: %v", err)
	}

	if len(jpegFiles) == 0 {
		return nil, fmt.Errorf("no JPEG files were generated")
	}

	return jpegFiles, nil
}
