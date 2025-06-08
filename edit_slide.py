#!/usr/bin/env python3
import uno
import sys
import os
from com.sun.star.connection import NoConnectException

def edit_slide_title(pptx_path, slide_index, new_title):
    """Edit slide title using UNO API"""
    try:
        # Connect to LibreOffice
        local_context = uno.getComponentContext()
        resolver = local_context.ServiceManager.createInstanceWithContext(
            "com.sun.star.bridge.UnoUrlResolver", local_context)
        
        # Connect to the running LibreOffice instance
        context = resolver.resolve("uno:socket,host=localhost,port=8100;urp;StarOffice.ComponentContext")
        desktop = context.ServiceManager.createInstanceWithContext(
            "com.sun.star.frame.Desktop", context)
        
        # Convert file path to file URL
        file_url = uno.systemPathToFileUrl(os.path.abspath(pptx_path))
        print(f"Loading document: {file_url}")
        
        # Load the presentation with headless properties
        from com.sun.star.beans import PropertyValue
        
        props = (
            PropertyValue("Hidden", 0, True, 0),
            PropertyValue("ReadOnly", 0, False, 0),
        )
        
        doc = desktop.loadComponentFromURL(file_url, "_blank", 0, props)
        
        # Get the slides
        slides = doc.getDrawPages()
        print(f"Found {slides.getCount()} slides")
        
        if slide_index >= slides.getCount():
            raise ValueError(f"Slide index {slide_index} out of range (0-{slides.getCount()-1})")
        
        # Get the specific slide
        slide = slides.getByIndex(slide_index)
        print(f"Editing slide {slide_index}")
        
        # Find and edit the title text
        # Usually the title is the first text shape on the slide
        for i in range(slide.getCount()):
            shape = slide.getByIndex(i)
            
            # Check if this shape has text and might be a title
            if hasattr(shape, 'getString') and hasattr(shape, 'setString'):
                current_text = shape.getString()
                print(f"Found text shape {i}: '{current_text}'")
                
                # If this looks like a title (often the first text shape or contains "Hello")
                if i == 0 or "Hello" in current_text:
                    print(f"Changing title from '{current_text}' to '{new_title}'")
                    shape.setString(new_title)
                    break
        
        # Save the document
        print("Saving document...")
        doc.store()
        
        # Close the document
        doc.close(True)
        print("Document saved and closed successfully")
        
    except NoConnectException:
        print("Error: Could not connect to LibreOffice. Make sure it's running with UNO socket.")
        sys.exit(1)
    except Exception as e:
        print(f"Error editing slide: {e}")
        sys.exit(1)

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: python3 edit_slide.py <pptx_path> <slide_index> <new_title>")
        sys.exit(1)
    
    pptx_path = sys.argv[1]
    slide_index = int(sys.argv[2])
    new_title = sys.argv[3]
    
    edit_slide_title(pptx_path, slide_index, new_title)
