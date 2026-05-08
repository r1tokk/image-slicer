package cmd

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	// Variables to store our flag values
	inputFile   string
	chunkHeight int
	prefix      string
	outputDir   string // <-- New variable for the target directory
)

// sliceCmd represents the slice command
var sliceCmd = &cobra.Command{
	Use:   "slice",
	Short: "Splits an image into vertical chunks",
	Long: `A CLI tool to split a long PNG image into smaller, page-sized chunks 
for easier reading or printing. 

Example usage:
  image-slicer slice -f "document.png" -H 800 -p "page-" -o "./output_folder"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Open the target image file
		file, err := os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer file.Close()

		// Decode the PNG
		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("error decoding image: %w", err)
		}

		// Calculate boundaries and chunks
		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		chunks := height / chunkHeight
		if height%chunkHeight != 0 {
			chunks++
		}

		// Interface required to crop standard images in Go
		type subImager interface {
			SubImage(r image.Rectangle) image.Image
		}

		sImg, ok := img.(subImager)
		if !ok {
			return fmt.Errorf("image format does not support cropping")
		}

		// Ensure the output directory exists (creates it if it doesn't)
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}

		// Loop through and slice the image
		for i := 0; i < chunks; i++ {
			startY := i * chunkHeight
			endY := (i + 1) * chunkHeight
			if endY > height {
				endY = height
			}

			// Define the crop rectangle
			rect := image.Rect(0, startY, width, endY)
			cropped := sImg.SubImage(rect)

			// Construct the full file path using filepath.Join
			filename := fmt.Sprintf("%s%d.png", prefix, i+1)
			fullPath := filepath.Join(outputDir, filename)

			out, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("error creating file %s: %w", fullPath, err)
			}

			if err := png.Encode(out, cropped); err != nil {
				out.Close()
				return fmt.Errorf("error encoding png: %w", err)
			}

			out.Close()
			fmt.Printf("Generated %s\n", fullPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sliceCmd)

	// Existing flags
	sliceCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Input PNG file to split (required)")
	sliceCmd.MarkFlagRequired("file")
	sliceCmd.Flags().IntVarP(&chunkHeight, "height", "H", 1200, "Height of each chunk in pixels")
	sliceCmd.Flags().StringVarP(&prefix, "prefix", "p", "chunk_", "Prefix for the output files")

	// New target directory flag
	sliceCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Target directory for saved chunks (default is current directory)")
}
