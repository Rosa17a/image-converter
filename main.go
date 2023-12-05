package main

import (
	"flag"
	"fmt"
	"github.com/davidbyttow/govips/v2/vips"
	"golang.org/x/image/webp"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Initialize vips
	vips.Startup(nil)
	defer vips.Shutdown()

	// Check for directory argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <finalFormat> <directory>")
		os.Exit(1)
	}

	inputFormat := flag.String("input", "", "Format of the input files")
	outputFormat := flag.String("output", "", "Format of the output files")
	// Parse the flags
	flag.Parse()
	dir := os.Args[3]

	log.Printf("Input format: %s", *inputFormat)
	log.Printf("Output format: %s", *outputFormat)

	// Read files in directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		os.Exit(1)
	}

	switch *inputFormat {
	case "jpeg":
		for _, file := range files {
			log.Printf("Processing %s", file.Name())
			if strings.HasSuffix(file.Name(), ".jpeg") {
				if *outputFormat == "png" {
					convertJpegToPng(dir, file.Name())
				}
				if *outputFormat == "webp" {
					toWebp(dir, file.Name())
				} else {
					panic("Unsupported format")
				}
			}
		}

	case "png":
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".png") {

				if *outputFormat == "jpeg" {
					convertPngToJpeg(dir, file.Name())
				} else {
					toWebp(dir, file.Name())
				}
			}
		}
	case "webp":
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".webp") {

				if *outputFormat == "png" {
					convertWebpToPng(dir, file.Name())
				} else {
					convertWebpToJpeg(dir, file.Name())
				}
			}
		}
	}
}

func convertPngToJpeg(dir string, name string) {
	// Open the PNG file
	filePath := dir + "/" + name
	pngFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(pngFile *os.File) {
		err := pngFile.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(pngFile)

	// Decode PNG
	img, err := png.Decode(pngFile)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Create a new JPEG file
	jpegFilename := strings.TrimSuffix(name, ".png") + ".jpeg"
	jpegFile, err := os.Create(dir + "/" + jpegFilename)
	if err != nil {
		fmt.Println("Error creating JPEG file:", err)
		return
	}
	defer func(jpegFile *os.File) {
		err := jpegFile.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(jpegFile)

	// Encode image to JPEG
	err = jpeg.Encode(jpegFile, img, nil)
	if err != nil {
		fmt.Println("Error encoding JPEG:", err)
	} else {
		fmt.Println("Converted:", jpegFilename)
	}
}

func toWebp(dir string, filename string) bool {
	// Construct the full path of the image
	fullPath := filepath.Join(dir, filename)

	// Load the image
	imageRef, err := vips.NewImageFromFile(fullPath)
	if err != nil {
		log.Printf("unable to open image %v: %v", fullPath, err)
		return true
	}

	params := vips.WebpExportParams{
		Quality:  100,
		Lossless: true,
	}

	// Convert to WebP using ExportNative
	webpBuf, _, err := imageRef.ExportWebp(&params)
	if err != nil {
		log.Printf("unable to convert image %v to webp: %v", fullPath, err)
		imageRef.Close()
		return true
	}
	imageRef.Close()

	// Construct the output filename
	outputFilename := strings.TrimSuffix(filename, filepath.Ext(filename)) + ".webp"
	outputPath := filepath.Join(dir, outputFilename)

	// Write the encoded buffer to a file
	err = os.WriteFile(outputPath, webpBuf, 0644)

	if err != nil {
		log.Printf("unable to write webp image %v: %v", outputPath, err)
		return true
	}
	return false
}

func convertJpegToPng(dir string, name string) {
	// Open the JPEG file
	filePath := dir + "/" + name
	jpegFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(jpegFile *os.File) {
		err := jpegFile.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(jpegFile)

	// Decode JPEG
	img, err := jpeg.Decode(jpegFile)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Create a new PNG file
	pngFilename := strings.TrimSuffix(name, ".jpeg") + ".png"
	pngFile, err := os.Create(dir + "/" + pngFilename)
	if err != nil {
		fmt.Println("Error creating PNG file:", err)
		return
	}
	defer func(jpegFile *os.File) {
		err := jpegFile.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(pngFile)

	// Encode image to JPEG
	err = png.Encode(pngFile, img)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
	} else {
		fmt.Println("Converted:", pngFilename)
	}
}

func convertWebpToJpeg(dir, filename string) {
	// Open the WEBP file
	filePath := dir + "/" + filename
	webpFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(webpFile *os.File) {
		err := webpFile.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(webpFile)

	// Decode WEBP
	img, err := webp.Decode(webpFile)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Create a new JPEG file
	jpegFilename := strings.TrimSuffix(filename, ".webp") + ".jpeg"
	jpegFile, err := os.Create(dir + "/" + jpegFilename)
	if err != nil {
		fmt.Println("Error creating JPEG file:", err)
		return
	}
	defer func(jpegFile *os.File) {
		err := jpegFile.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(jpegFile)

	// Encode image to JPEG
	err = jpeg.Encode(jpegFile, img, nil)
	if err != nil {
		fmt.Println("Error encoding JPEG:", err)
	} else {
		fmt.Println("Converted:", jpegFilename)
	}
}

func convertWebpToPng(dir, filename string) {
	// Open the WEBP file
	filePath := dir + "/" + filename
	webpFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer func(webpFile *os.File) {
		err := webpFile.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(webpFile)

	// Decode WEBP
	img, err := webp.Decode(webpFile)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Create a new PNG file
	pngFilename := strings.TrimSuffix(filename, ".webp") + ".png"
	pngFile, err := os.Create(dir + "/" + pngFilename)
	if err != nil {
		fmt.Println("Error creating PNG file:", err)
		return
	}
	defer func(jpegFile *os.File) {
		err := jpegFile.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(pngFile)

	// Encode image to JPEG
	err = png.Encode(pngFile, img)
	if err != nil {
		fmt.Println("Error encoding PNG:", err)
	} else {
		fmt.Println("Converted:", pngFilename)
	}
}
