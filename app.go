package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

type FileDetails struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Size     int64  `json:"size"` // bytes
	Content  string `json:"content"`
}

type ReturnResult struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Html     string `json:"html"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SelectAndReadFile() (*FileDetails, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select a file",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Markdown Files",
				Pattern:     "*.md",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	// User cancelled the dialog
	if path == "" {
		return nil, errors.New("file selection cancelled")
	}

	// Get file metadata
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Read file content
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &FileDetails{
		Filename: filepath.Base(path),
		Path:     path,
		Size:     info.Size(),
		Content:  string(content), // Cast []byte to string for text files
	}, nil
}

func (a *App) ConvertToMd(file *FileDetails) (*ReturnResult, error) {
	// Placeholder for conversion logic
	file, err := a.SelectAndReadFile()
	if err != nil {
		return nil, err
	}

	htmlContent := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Typographer, extension.Linkify), // Add GitHub Flavored Markdown extension
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Automatically generate IDs for headings
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // Render hard line breaks as <br>
			html.WithXHTML(),     // Use XHTML-style self-closing tags
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // Allow raw HTML in Markdown (be cautious with this)
		),
	)
	var buf bytes.Buffer
	if err := htmlContent.Convert([]byte(file.Content), &buf); err != nil {
		return nil, err
	}

	runtime.WindowSetTitle(a.ctx, file.Path+" - Margo")

	return &ReturnResult{
		Filename: file.Filename,
		Path:     file.Path,
		Html:     buf.String(), // Use the converted HTML content

	}, nil
}
