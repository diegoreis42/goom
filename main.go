package main

import (
    "bufio"
    "fmt"
    "image"
    "os"
    "os/exec"
    "runtime"
    "time"

    "gocv.io/x/gocv"
)

const (
    asciiChars = " .:-=+*#%@"
)

func clearScreen() {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("cmd", "/c", "cls")
    default:
        cmd = exec.Command("clear")
    }
    cmd.Stdout = os.Stdout
    _ = cmd.Run()
}

func mapPixelToAscii(pixelValue uint8) byte {
    idx := int(pixelValue) * (len(asciiChars) - 1) / 255
    return asciiChars[idx]
}

func main() {
    // Open webcam
    webcam, err := gocv.VideoCaptureDevice(0)
    if err != nil {
        fmt.Println("Error opening webcam:", err)
        os.Exit(1)
    }
    defer webcam.Close()

    if !webcam.IsOpened() {
        fmt.Println("Error: Could not open webcam.")
        os.Exit(1)
    }

    // Create a Mat to store frames
    img := gocv.NewMat()
    defer img.Close()

    // Get terminal size
    width, height := 80, 24

    clearScreen()
    fmt.Println("Press Ctrl+C to exit...")

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

        // Use buffered output for better performance
        writer := bufio.NewWriter(os.Stdout)

        // Convert image to ASCII art
        for y := 0; y < img.Rows(); y++ {
            for x := 0; x < img.Cols(); x++ {
                pixel := img.GetUCharAt(y, x)
                writer.WriteByte(mapPixelToAscii(pixel))
            }
            writer.WriteByte('\n')
        }

        writer.Flush()

        // Add a small delay to limit frame rate and prevent high CPU usage
        time.Sleep(50 * time.Millisecond)

        // Move cursor to top left to overwrite the previous frame
        fmt.Print("\033[H")
    }
}



