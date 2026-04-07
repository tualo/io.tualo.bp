package paper

import (
	"image"
	"math"

	"gocv.io/x/gocv"
)

func (me *Paper) GetCornerPoints(contour gocv.PointVector) map[string]image.Point {

	topLeftCorner := image.Point{}
	topLeftCornerDist := 0.0

	topRightCorner := image.Point{}
	topRightCornerDist := 0.0

	bottomLeftCorner := image.Point{}
	bottomLeftCornerDist := 0.0

	bottomRightCorner := image.Point{}
	bottomRightCornerDist := 0.0

	if contour.Size() > 0 {
		rect := gocv.MinAreaRect(contour)
		center := rect.Center
		// fmt.Println("center: ", center)
		for i := 0; i < contour.Size(); i++ {
			point := contour.At(i)

			//dist := gocv.Norm(point, center, gocv.NormL2)
			dist := math.Hypot(math.Abs(float64(point.X-center.X)), math.Abs(float64(point.Y-center.Y)))

			if point.X < center.X && point.Y < center.Y {
				// top left
				if dist > topLeftCornerDist {
					topLeftCorner = point
					topLeftCornerDist = dist
				}
			} else if point.X > center.X && point.Y < center.Y {
				// top right
				if dist > topRightCornerDist {
					topRightCorner = point
					topRightCornerDist = dist
				}
			} else if point.X < center.X && point.Y > center.Y {
				// bottom left
				if dist > bottomLeftCornerDist {
					bottomLeftCorner = point
					bottomLeftCornerDist = dist
				}
			} else if point.X > center.X && point.Y > center.Y {
				// bottom right
				if dist > bottomRightCornerDist {
					bottomRightCorner = point
					bottomRightCornerDist = dist
				}
			}

		}
	}

	return map[string]image.Point{
		"topLeftCorner":     topLeftCorner,
		"topRightCorner":    topRightCorner,
		"bottomLeftCorner":  bottomLeftCorner,
		"bottomRightCorner": bottomRightCorner,
	}
}
