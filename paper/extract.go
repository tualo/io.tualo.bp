package paper

import (
	// "fmt"
	"image"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
	//"image/color"
)

func (me *Paper) Extract(
	img *cv.Mat,
	cornerPoints map[string]image.Point,
) (cv.Mat, cv.Mat) {

	warpedDst := cv.NewMat("warpedDst")
	var M cv.Mat
	var invM cv.Mat

	topLeftCorner := cornerPoints["topLeftCorner"]
	topRightCorner := cornerPoints["topRightCorner"]
	bottomLeftCorner := cornerPoints["bottomLeftCorner"]
	bottomRightCorner := cornerPoints["bottomRightCorner"]

	resultWidth := topRightCorner.X - topLeftCorner.X
	resultHeight := bottomRightCorner.Y - topRightCorner.Y

	dsize := image.Point{resultWidth, resultHeight}
	if topLeftCorner != (image.Point{}) && topRightCorner != (image.Point{}) && bottomLeftCorner != (image.Point{}) && bottomRightCorner != (image.Point{}) {
		newImg := []image.Point{
			{0, 0},
			{0, resultHeight},
			{resultWidth, resultHeight},
			{resultWidth, 0},
		}
		origImg := []image.Point{
			topLeftCorner,     // top-left
			bottomLeftCorner,  // bottom-left
			bottomRightCorner, // bottom-right
			topRightCorner,    // top-right
		}
		origV := gocv.NewPointVectorFromPoints(origImg)
		newV := gocv.NewPointVectorFromPoints(newImg)

		M = cv.NewFromMat("M", gocv.GetPerspectiveTransform(origV, newV))
		invM = cv.NewFromMat("invM", gocv.GetPerspectiveTransform(newV, origV))
		gocv.WarpPerspective(img.Get(), warpedDst.GetPointer(), M.Get(), dsize)

		origV.Close()
		newV.Close()
		M.Close()
	}
	return warpedDst, invM
}
