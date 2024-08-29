package main

import (
    "bufio"
    "fmt"
    "image"
    "os"
    "os/exec"
    "time"

    "gocv.io/x/gocv"
)

const (
  asciiChars = " ._,:;ox=%#@"
)

func clearScreen() {
    cmd := exec.Command("clear")
    cmd.Stdout = os.Stdout
    _ = cmd.Run()
}

func mapPixelToAscii(pixelValue uint8) byte {
    idx := int(pixelValue) * (len(asciiChars) - 1) / 255
    return asciiChars[idx]
}

func main() {
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

    img := gocv.NewMat()
    defer img.Close()

    width, height := 100, 24

    clearScreen()

    for {
        if ok := webcam.Read(&img); !ok || img.Empty() {
            fmt.Println("Error reading image from webcam")
            continue
        }

        gocv.Resize(img, &img, image.Point{width, height}, 0, 0, gocv.InterpolationDefault)

        gocv.CvtColor(img, &img, gocv.ColorBGRToGray)

        writer := bufio.NewWriter(os.Stdout)

        for y := 0; y < img.Rows(); y++ {
            for x := 0; x < img.Cols(); x++ {
                pixel := img.GetUCharAt(y, x)
                writer.WriteByte(mapPixelToAscii(pixel))
            }
            writer.WriteByte('\n')
        }

        writer.Flush()

        time.Sleep(25 * time.Millisecond)

        fmt.Print("\033[H")
    }
}

