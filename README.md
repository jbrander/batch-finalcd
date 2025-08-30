# Batch FinalCD Frontend

A modern desktop frontend for FinalCD audio downsampling software, featuring batch processing capabilities.

## Overview

FinalCD is professional audio downsampling software developed by Sonic Illusions. This frontend application provides an enhanced user interface with batch processing support, allowing you to convert multiple WAV files simultaneously instead of processing them one at a time.

## Requirements

- Windows operating system
- FinalCD executable (`finalcd.exe`) must be in the same directory as this application
- WAV audio files for processing

## Installation

1. Download the latest release from the GitHub releases page
2. Extract the executable to a directory containing `finalcd.exe`
3. Run `batch-finalcd-frontend-windows-amd64.exe`

## Usage

1. **Add Input Files**: Click "Add" to select multiple WAV files using the file dialog
2. **Set Output Directory**: The output directory defaults to a "finalcd_output" folder next to your input files, or click "..." to choose a custom location
3. **Configure Processing Options**:
   - Select filter type (Sharp, Gentle, or Alternative)
   - Choose dither setting (None or TPDF)
   - Enable 32-bit output if needed
   - Adjust gain value if required
4. **Convert**: Click "Convert!" to begin batch processing
5. **Monitor Progress**: Watch the output log for real-time processing status

## File Structure

```
batch-finalcd-frontend/
├── frontend/           # Vite-based frontend application
│   ├── src/
│   │   ├── main.js    # Main application logic
│   │   └── style.css  # Custom styles
│   ├── index.html     # Main UI template
│   └── package.json   # Frontend dependencies
├── app.go             # Go backend with FinalCD integration
├── main.go            # Wails application entry point
├── wails.json         # Wails configuration
└── go.mod             # Go module dependencies
```

## Technical Details

This application is built using:
- **Wails v2**: Go-based framework for building desktop applications with web frontends
- **Go 1.24**: Backend runtime for file operations and FinalCD integration
- **Vite**: Fast frontend build tool
- **Tailwind CSS**: Utility-first CSS framework for styling
- **Vanilla JavaScript**: No heavy frontend frameworks, keeping the application lightweight

## Building from Source

### Prerequisites
- Go 1.24+ 
- Node.js 22+
- Wails CLI v2.10.2

### Build Steps
```bash
# Install Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.2

# Install frontend dependencies
cd frontend && npm install && cd ..

# Build for development
wails dev

# Build for production (Windows)
wails build -platform windows/amd64
```

## About FinalCD

FinalCD is professional audio downsampling software created by Sonic Illusions. It provides high-quality sample rate conversion with various filtering and dithering options. The original software can be found at http://www.sonicillusions.co.uk/finalcd.htm.

This frontend application enhances the original by adding batch processing capabilities, making it easier to process large numbers of audio files efficiently.

## License

This frontend application is released into the public domain. Please refer to the original FinalCD documentation for licensing information regarding the core audio processing engine.
