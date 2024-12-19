package splitroi

import (
	"fmt"
	"image"
	"log"

	"gocv.io/x/gocv"
)

func extractItem(
	img gocv.Mat,
	roi ROI,
	pixelScaleX float64,
	pixelScaleY float64,
	// nextStep func(gocv.Mat, splitroi.ROI, float64, float64),
) {

	for y := 0; y < roi.ItemCountY; y++ {
		for x := 0; x < roi.ItemCountX; x++ {
			log.Println(x, y)
			X := int(roi.X*pixelScaleX*0 + roi.Width*pixelScaleX*float64(x))
			Y := int(roi.Y*pixelScaleY*0 + roi.Height*pixelScaleY*float64(y))
			W := int(roi.Width * pixelScaleX)
			H := int(roi.Height * pixelScaleY)
			rect := image.Rect(X, Y, X+W, Y+H)
			croppedMat := img.Region(rect)
			if !croppedMat.Empty() {
				// nextStep(croppedMat, roi, pixelScaleX, pixelScaleY)
				fileName := fmt.Sprintf("%d.%d.jpg", x, y)
				gocv.IMWrite(fileName, croppedMat)
				circles(croppedMat, pixelScaleX, pixelScaleY)
				fileNameC := fmt.Sprintf("%d.%d.c.jpg", x, y)
				gocv.IMWrite(fileNameC, croppedMat)
			}
		}
	}
}
