package freezer

import (
	"image"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
	ppr "tualo.de/deep-test/paper"
)

type Freezer struct {
	mog    gocv.BackgroundSubtractorMOG2
	frozen bool
}

const MinimumArea = 8000

var paper ppr.Paper = ppr.Paper{}

func (me *Freezer) Init() {
	me.mog = gocv.NewBackgroundSubtractorMOG2WithParams(50, 200, true)

	paper.Init()
}

/*
this function returns true if there is no movment detected anymore
*/
func (me *Freezer) Check(img *cv.Mat) bool {
	//	log.Println("Freezer start", cv.MatCounter)

	me.frozen = true

	factor := 12
	scaled := cv.NewMat("scaled freezer")

	imgDelta := gocv.NewMat()
	imgThresh := gocv.NewMat()
	gocv.Resize(img.Get(), scaled.GetPointer(), image.Point{img.Cols() / factor, img.Rows() / factor}, 0, 0, gocv.InterpolationArea) // faster detection

	me.mog.Apply(scaled.Get(), &imgDelta)
	scaled.Close()
	// remaining cleanup of the image to use for finding contours.
	// first use threshold
	gocv.Threshold(imgDelta, &imgThresh, 120, 255, gocv.ThresholdBinary)
	// then dilate
	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
	gocv.Dilate(imgThresh, &imgThresh, kernel)
	kernel.Close()

	// now find contours
	contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	for i := 0; i < contours.Size(); i++ {
		area := gocv.ContourArea(contours.At(i))

		if area < MinimumArea {
			continue
		}
		me.frozen = false
	}

	contours.Close()
	imgThresh.Close()
	imgDelta.Close()

	return me.frozen
}
