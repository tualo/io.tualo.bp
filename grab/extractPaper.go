package grab

import (
	// "fmt"
	"image"
	"gocv.io/x/gocv"
  //"image/color"
)

func extractPaper(img gocv.Mat, maxContour gocv.PointVector, resultWidth int, resultHeight int, cornerPoints map[string]image.Point) (gocv.Mat, gocv.Mat) {
	topLeftCorner := cornerPoints["topLeftCorner"]
	topRightCorner := cornerPoints["topRightCorner"]
	bottomLeftCorner := cornerPoints["bottomLeftCorner"]
	bottomRightCorner := cornerPoints["bottomRightCorner"]
	warpedDst := gocv.NewMat()
	M := gocv.NewMat()
	invM := gocv.NewMat()
  dsize := image.Point{resultWidth, resultHeight}
	if topLeftCorner != (image.Point{}) && topRightCorner != (image.Point{}) && bottomLeftCorner != (image.Point{}) && bottomRightCorner != (image.Point{}) {
    newImg := []image.Point{
      image.Point{0, 0},
      image.Point{0, resultHeight},
      image.Point{resultWidth, resultHeight},
      image.Point{resultWidth, 0},
    }
    origImg := []image.Point{
      topLeftCorner, // top-left
      bottomLeftCorner, // bottom-left
      bottomRightCorner, // bottom-right
      topRightCorner,  // top-right
    }
    origV := gocv.NewPointVectorFromPoints(origImg)
    newV := gocv.NewPointVectorFromPoints(newImg)

    M.Close()
		M = gocv.GetPerspectiveTransform( origV  , newV)
    invM = gocv.GetPerspectiveTransform(  newV, origV)
		gocv.WarpPerspective(img, &warpedDst, M, dsize)

    //WarpPerspectiveWithParams(src Mat, dst *Mat, m Mat, sz image.Point, flags InterpolationFlags, borderType BorderType, borderValue color.RGBA)
    //gocv.BorderDefault
    //gocv.WarpInverseMap

    origV.Close()
    newV.Close()
    M.Close()
	}
	return warpedDst, invM
}
