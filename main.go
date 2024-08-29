package main

import (
    "fmt"
    "image"
    "os"
    "gocv.io/x/gocv"
)

const (
    asciiChars = " .:-=+*#123%@"
    width      =  80 // Width of the ASCII art in characters
    height     = 24  // Height of the ASCII art in characters
)

func main() {
    // Open webcam
    webcam, err := gocv.VideoCaptureDevice(0)
    if err != nil {
        fmt.Println("Error opening webcam:", err)
        os.Exit(1)
    }
    defer webcam.Close()

    // Create a window to display the original image
    window := gocv.NewWindow("Webcam")
    defer window.Close()

    // Create a Mat to store frames
    img := gocv.NewMat()
    defer img.Close()
    // Main loop
    for {
        if ok := webcam.Read(&img); !ok || img.Empty() {
            fmt.Println("Error reading image from webcam")
            continue
        }

        // Resize the image to fit the terminal
        gocv.Resize(img, &img, image.Point{width, height}, 0, 0, gocv.InterpolationDefault)

        // Convert the image to grayscale
        gocv.CvtColor(img, &img, gocv.ColorBGRToGray)

        // Clear the terminal
        fmt.Print("\033[H\033[2J")

        // Convert image to ASCII art
        for y := 0; y < img.Rows(); y++ {
            for x := 0; x < img.Cols(); x++ {
                pixel := img.GetUCharAt(y, x)
                idx := int(pixel) * len(asciiChars) / 256
                fmt.Printf("%c", asciiChars[idx])
            }
            fmt.Println()
        }
    }
}

