# image-slicer

A CLI tool written in Go that splits long PNG images into smaller, page-sized vertical chunks. This is especially useful for breaking down long documents, infographics, or webcomics for easier reading and printing.

## Features

* **Vertical Slicing:** Automatically calculates and splits an image into multiple chunks based on a target pixel height.
* **Customizable Output:** Easily specify the output directory and a custom prefix for your generated chunks.

## Installation

The project includes a `build.sh` script that compiles the Go binary and installs it globally on your system.

1. Clone the repository:
   ```bash
   git clone git@github.com:r1tokk/image-slicer.git
   cd image-slicer
2. Run the build script (you may need sudo since it copies the executable to /usr/local/bin):
    ```bash 
    chmod +x build.sh
    sudo ./build.sh

## Usage

The primary command to split images is `slice`:

    image-slicer slice -f <input_file.png> [flags]    

### Basic Usage
    image-slicer slice -f "document.png"

### Custom Height and Prefix:
    image-slicer slice -f "webcomic.png" -H 800 -p "page-"

### Specify an Output Directory:
    image-slicer slice -f "infographic.png" -o "./output_folder"

### Usage with Freeze (Slicing Code Screenshots)

If you use [freeze](https://github.com/charmbracelet/freeze) to generate beautiful, long screenshots of your code, you can easily slice them into paginated chunks. 

Because `freeze` determines its output format based on the file extension (defaulting to SVG if piped), the most reliable way to slice its output is by chaining commands to create, process, and delete a temporary `.png` file in one go:

```bash
freeze main.go -o temp.png && image-slicer slice -f temp.png -H 800 -p "code-page-" -o "./code_chunks" && rm temp.png
