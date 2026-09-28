package cmd

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	convertInputFile  string
	convertOutputFile string
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts an image from one format to another",
	Long: `A CLI tool to convert images between formats (e.g., PNG to JPEG). 
The output format is determined by the output file's extension.

Example usage:
  image-slicer convert -f "document.png" -o "document.jpg"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Open the input file
		file, err := os.Open(convertInputFile)
		if err != nil {
			return fmt.Errorf("%s %w", color.RedString("error opening file:"), err)
		}
		defer file.Close()

		// Decode the image
		// image.Decode relies on imported image packages (image/png, image/jpeg)
		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("%s %w", color.RedString("error decoding image:"), err)
		}

		// Create the output file
		out, err := os.Create(convertOutputFile)
		if err != nil {
			return fmt.Errorf("%s %w", color.RedString("error creating file %s:", convertOutputFile), err)
		}
		defer out.Close()

		// Determine format from the output file extension
		ext := strings.ToLower(filepath.Ext(convertOutputFile))
		switch ext {
		case ".png":
			err = png.Encode(out, img)
		case ".jpg", ".jpeg":
			err = jpeg.Encode(out, img, &jpeg.Options{Quality: 100})
		default:
			return fmt.Errorf(color.RedString("unsupported output format: %s", ext))
		}

		if err != nil {
			return fmt.Errorf("%s %w", color.RedString("error encoding image:"), err)
		}

		color.Cyan("Successfully converted %s to %s\n", convertInputFile, convertOutputFile)
		return nil
	},
}

func init() {
	// Register the command with the root command
	rootCmd.AddCommand(convertCmd)

	// Define command options (flags)
	convertCmd.Flags().StringVarP(&convertInputFile, "file", "f", "", "Input image file to convert (required)")
	convertCmd.MarkFlagRequired("file")

	convertCmd.Flags().StringVarP(&convertOutputFile, "output", "o", "", "Output image file path with desired extension (e.g., output.jpg) (required)")
	convertCmd.MarkFlagRequired("output")
}
