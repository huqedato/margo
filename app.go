package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed wails.json
var wailsConfig []byte

// App struct
type App struct {
	ctx             context.Context
	initialFilePath string
}

type AppInfo struct {
	Info struct {
		ProductName    string `json:"productName"`
		ProductVersion string `json:"productVersion"`
		CompanyName    string `json:"companyName"`
		Copyright      string `json:"copyright"`
		Description    string `json:"description"`
		License        string `json:"license"`
	} `json:"info"`
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
	app := &App{}
	if len(os.Args) > 1 {
		app.initialFilePath = os.Args[1]
	}
	return app
}

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

func (a *App) GetStartupFile() (*ReturnResult, error) {
	if a.initialFilePath == "" {
		return nil, nil
	}

	// Get file metadata
	info, err := os.Stat(a.initialFilePath)
	if err != nil {
		return nil, err
	}

	// Read file content
	content, err := os.ReadFile(a.initialFilePath)
	if err != nil {
		return nil, err
	}

	fileDetails := &FileDetails{
		Filename: filepath.Base(a.initialFilePath),
		Path:     a.initialFilePath,
		Size:     info.Size(),
		Content:  string(content),
	}

	return a.ConvertToMd(fileDetails)
}

func (a *App) OpenFileFromFrontend(file *FileDetails) (*ReturnResult, error) {
	file, err := a.SelectAndReadFile()
	if err != nil {
		return nil, err
	}
	return a.ConvertToMd(file)
}

func (a *App) ConvertToMd(file *FileDetails) (*ReturnResult, error) {

	htmlContent := goldmark.New(
		goldmark.WithExtensions(extension.GFM, extension.Footnote, extension.Typographer, extension.Linkify,
			highlighting.NewHighlighting(
				highlighting.WithStyle("onedark"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(false),
					chromahtml.WithCustomCSS(map[chroma.TokenType]string{
						chroma.Background: "background-color: transparent;", // Remove background color from code blocks
						chroma.PreWrapper: "padding: 1em; border-radius: 5px;",
					}),
				),
			),
		), // Add GitHub Flavored Markdown extension
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
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
		Html:     buf.String(),
	}, nil
}

func (a *App) GetAppInfo() (*AppInfo, error) {
	var info AppInfo
	err := json.Unmarshal(wailsConfig, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
