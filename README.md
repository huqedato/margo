# Margo

Margo is an opinionated Markdown viewer built with Wails, Go, and Vue.

It is a single-window reader for local `.md` files with a deliberately small surface area: open a file, render it quickly, search inside it, zoom the document, and get out of the way. No Electron, no fuss, just a simple, distraction-free, and efficient Markdown viewer for desktop.

## What It Does

- Opens Markdown files from a native file picker.
- Accepts a file path as the first CLI argument and loads it on startup.
- Supports drag and drop for local files.
- Renders Markdown with Goldmark using:
	- GitHub Flavored Markdown
	- Footnotes
	- Typographer
	- Linkify
	- Auto heading IDs
	- Syntax highlighting via Chroma using the `onedark` style
- Provides in-document search using `mark.js`.
- Supports simple zoom controls.
- Sends clicked links to the system browser.

## What It Does Not Do

This repository currently implements a focused viewer, not a full Markdown workspace.

- No editing.
- No tabbed browsing or document history.
- No live file watching or auto-refresh when the source file changes.
- No print/export flow.
- No in-app navigation for local Markdown links; rendered anchors are handed off to the external browser.


## Building latest

Build the desktop app for Windows x64 from the project root:

```bash
wails build
```

The Wails packaging assets and platform-specific metadata live under `build/`.

## Usage

Start Margo normally and open a file from the UI, or launch it with a Markdown path:

```bash
margo path/to/file.md
```

On Windows that will typically be:

```powershell
.\margo.exe C:\path\to\file.md
```

Drag and drop a Markdown file onto the open app window to load it.


## Status

The current Margo implementation is minimal and I want to keep it that way. Actually this is the intended scope of the project. Future improvements will be focused on bug fixes and small quality-of-life adjustments, and less on adding new features.

Margo is currently built for Windows, but the architecture is cross-platform and should be portable to macOS and Linux with minimal adjustments. 

## Contributing
Contributions are welcome! If you have an idea for a new feature or improvement, please open an issue or submit a pull request. But mind to keep this minimalistic and usable. 

## Not Vibecoded

This project is not vibecoded.

## License

This software is released under **[AGPL-3.0-or-later](https://www.gnu.org/licenses/agpl-3.0.html)**.