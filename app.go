package main

import (
	"bytes"
	"context"
	"strings"

	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

type FinalCdResult struct {
	InputFileName  string
	OutputFilePath string
	Success        bool
	Stdout         string
	Err            string
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ResizeWindow(width int, height int) {
	runtime.WindowSetSize(a.ctx, width, height)
}

func (a *App) SelectOutputDirectory() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:                "Select Output Directory",
		CanCreateDirectories: true,
		ShowHiddenFiles:      true,
		ResolvesAliases:      true,
	})
	if err != nil {
		return ""
	}
	return dir
}

func (a *App) DefaultOutputDirectory(inputFilePath string) string {
	dir := filepath.Dir(inputFilePath)

	return filepath.Join(dir, "finalcd_output")
}

func (a *App) SelectInputFiles() []string {
	dialogOptions := runtime.OpenDialogOptions{
		Title:           "Select Input WAV Files",
		ShowHiddenFiles: true,
		ResolvesAliases: true,
		Filters: []runtime.FileFilter{
			{DisplayName: "Wave Files (*.wav)", Pattern: "*.wav"},
		},
	}
	selectedFiles, err := runtime.OpenMultipleFilesDialog(a.ctx, dialogOptions)
	if err != nil {
		return nil
	}
	return selectedFiles
}

func (a *App) NativeAlert(title string, message string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})
}

func (a *App) RunFinalCD(inputFilePath string, outputDirectory string, options []string) FinalCdResult {
	inputFileName := filepath.Base(inputFilePath)
	outputFilePath := filepath.Join(outputDirectory, inputFileName)
	result := FinalCdResult{
		InputFileName:  inputFileName,
		OutputFilePath: outputFilePath,
		Success:        false,
		Stdout:         "",
		Err:            "",
	}

	if inputFilePath == "" {
		result.Err = "Input file path is empty"
		return result
	} else if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		result.Err = "Input file does not exist"
		return result
	} else if ext := filepath.Ext(inputFilePath); strings.ToLower(ext) != ".wav" {
		result.Err = "Input file is not a WAV file"
		return result
	}

	if outputDirectory == "" {
		result.Err = "No output directory selected"
		return result
	} else {
		// Create the output directory if it doesn't exist
		err := os.MkdirAll(outputDirectory, os.ModePerm)
		if err != nil {
			result.Err = fmt.Sprintf("Failed to create output directory: %v", err)
			return result
		}
	}

	// get the path of the current executable
	currentExePath, err := os.Executable()
	if err != nil {
		result.Err = fmt.Sprintf("Failed to get current executable path: %v", err)
		return result
	}
	// construct the path to finalcd.exe
	finalcdExePath := filepath.Join(filepath.Dir(currentExePath), "finalcd.exe")
	if _, err := os.Stat(finalcdExePath); os.IsNotExist(err) {
		result.Err = "finalcd.exe not found"
		return result
	}

	// Capture stdout
	var stdout bytes.Buffer
	cliArgs := buildCliArgs(inputFilePath, outputFilePath, options)

	// cmd := exec.Command(finalcdExePath, cliArgs...)
	// cmd.Stdout = &stdout
	// err = cmd.Run()

	err = runCommand(finalcdExePath, cliArgs, &stdout)

	result.Stdout += normalizeCarriageReturns(stdout)

	if err != nil {
		result.Err = fmt.Sprintf("Failed to run command: %v\n", err)
		return result
	}

	result.Success = true
	return result
}

func buildCliArgs(inputFilePath string, outputFilePath string, options []string) []string {
	var args []string
	for _, option := range options {
		if option != "" {
			args = append(args, option)
		}
	}
	args = append(args, inputFilePath, outputFilePath)
	return args
}

func normalizeCarriageReturns(buffer bytes.Buffer) string {
	input := strings.ReplaceAll(buffer.String(), "\r\n", "\n")
	var result strings.Builder
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// If the line contains carriage returns, keep only the last segment
		if idx := strings.LastIndex(line, "\r"); idx != -1 {
			line = line[idx+1:]
		}
		result.WriteString(line)
		result.WriteString("\n")
	}
	return result.String()
}
