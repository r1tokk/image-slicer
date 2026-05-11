package cmd

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	inputFile   string
	chunkHeight int
	prefix      string
	outputDir   string
)

var sliceCmd = &cobra.Command{
	Use:   "slice",
	Short: "Splits an image into vertical chunks",
	Long: `A CLI tool to split a long PNG image into smaller, page-sized chunks 
for easier reading or printing. 

Example usage:
  image-slicer slice -f "document.png" -H 800 -p "page-" -o "./output_folder"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("%s %w", color.RedString("error opening file:"), err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("%s %w", color.RedString("error opening file:"), err)
		}

		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		chunks := height / chunkHeight
		if height%chunkHeight != 0 {
			chunks++
		}

		type subImager interface {
			SubImage(r image.Rectangle) image.Image
		}

		sImg, ok := img.(subImager)
		if !ok {
			return fmt.Errorf(color.RedString("image format does not support cropping"))
		}

		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("%s %w", color.RedString("error creating output directory:"), err)
		}

		for i := 0; i < chunks; i++ {
			startY := i * chunkHeight
			endY := (i + 1) * chunkHeight
			if endY > height {
				endY = height
			}

			rect := image.Rect(0, startY, width, endY)
			cropped := sImg.SubImage(rect)

			filename := fmt.Sprintf("%s%d.png", prefix, i+1)
			fullPath := filepath.Join(outputDir, filename)

			out, err := os.Create(fullPath)
			if err != nil {
				return fmt.Errorf("%s %w", color.RedString("error creating file %s:", fullPath), err)
			}

			if err := png.Encode(out, cropped); err != nil {
				out.Close()
				return fmt.Errorf("%s %w", color.RedString("error encoding png:"), err)
			}

			out.Close()
			color.Cyan("Generated %s\n", fullPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sliceCmd)

	sliceCmd.Flags().StringVarP(&inputFile, "file", "f", "", "Input PNG file to split (required)")
	sliceCmd.MarkFlagRequired("file")
	sliceCmd.Flags().IntVarP(&chunkHeight, "height", "H", 1200, "Height of each chunk in pixels")
	sliceCmd.Flags().StringVarP(&prefix, "prefix", "p", "chunk_", "Prefix for the output files")

	sliceCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Target directory for saved chunks (default is current directory)")
}
