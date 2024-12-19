package main

import (
	"log"
	"os"

	"gocv.io/x/gocv"

	splitroi "io.tualo.bp/splitroi"
)

func main() {
	filename := os.Args[1]
	log.Println("filename", filename)
	img := gocv.IMRead(filename, gocv.IMReadColor)

	page_width := 210
	page_height := 400

	scaleX := float64(img.Cols()) / float64(page_width)
	scaleY := float64(img.Rows()) / float64(page_height)

	rois := splitroi.ROIS{
		Roi: []splitroi.ROI{
			{

				X:          81,
				Y:          54,
				Height:     21,
				Width:      21,
				ItemCountX: 1,
				ItemCountY: 15,
				XCap:       0.00000,
				YCap:       0.02000,
			},
		},
	}

	log.Println("Extract", scaleX, scaleY)
	splitroi.Extract(img, rois, scaleX, scaleY)
}
