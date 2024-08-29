package main

import (
    "bufio"
    "fmt"
    "image"
    "os"
    "os/exec"
    "runtime"
    "sync"
    "time"

    "gocv.io/x/gocv"
)

const (
    asciiChars = ".'`^,:;Il!i><~+_-?][}{1)(|/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

    width      = 80 // Adjust width for higher resolution
    height     = 36  // Adjust height for higher resolution
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

func processRow(img gocv.Mat, y int, output chan<- string, wg *sync.WaitGroup) {
    defer wg.Done()
    row := make([]byte, img.Cols())
    for x := 0; x < img.Cols(); x++ {
        pixel := img.GetUCharAt(y, x)
        row[x] = mapPixelToAscii(pixel)
    }
    output <- string(row)
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

    clearScreen()
    fmt.Println("Press Ctrl+C to exit...")

    for {
        if ok := webcam.Read(&img); !ok || img.Empty() {
            fmt.Println("Error reading image from webcam")
            continue
        }

        // Resize the image to fit the terminal
        gocv.Resize(img, &img, image.Point{width, height}, 0, 0, gocv.InterpolationArea)

        // Convert the image to grayscale
        gocv.CvtColor(img, &img, gocv.ColorBGRToGray)

        output := make(chan string, img.Rows())
        var wg sync.WaitGroup

        // Launch a goroutine for each row of the image
        for y := 0; y < img.Rows(); y++ {
            wg.Add(1)
            go processRow(img, y, output, &wg)
        }

        // Wait for all goroutines to finish
        go func() {
            wg.Wait()
            close(output)
        }()

        // Use buffered output for better performance
        writer := bufio.NewWriter(os.Stdout)

        // Write the ASCII art row by row
        for row := range output {
            writer.WriteString(row + "\n")
        }

        writer.Flush()

        // Add a small delay to limit frame rate and prevent high CPU usage
        time.Sleep(33 * time.Millisecond)

        // Move cursor to top left to overwrite the previous frame
        fmt.Print("\033[H")
    }
}

