package splitroi

import (
	"fmt"
	"image"
	"log"

	"gocv.io/x/gocv"
)

func Extract(
	img gocv.Mat,
	rois ROIS,
	pixelScaleX float64,
	pixelScaleY float64,
	// nextStep func(gocv.Mat, ROI, float64, float64, func(gocv.Mat, ROI, float64, float64)),
	// nextStep2 func(gocv.Mat, ROI, float64, float64),
) {
	log.Println("Extract")
	for useRoi := 0; useRoi < len(rois.Roi); useRoi++ {
		X := int(rois.Roi[useRoi].X * pixelScaleX)
		Y := int(rois.Roi[useRoi].Y * pixelScaleY)
		W := int(rois.Roi[useRoi].Width * pixelScaleX * float64(rois.Roi[useRoi].ItemCountX))
		H := int(rois.Roi[useRoi].Height * pixelScaleY * float64(rois.Roi[useRoi].ItemCountY))
		rect := image.Rect(X, Y, X+W, Y+H)
		log.Println("Extract", rect)
		croppedMat := img.Region(rect)
		log.Println(useRoi)
		if !croppedMat.Empty() {
			fileName := fmt.Sprintf("ROI.%d.jpg", useRoi)
			gocv.IMWrite(fileName, croppedMat)
			extractItem(croppedMat, rois.Roi[useRoi], pixelScaleX, pixelScaleY)
		}
	}
}
