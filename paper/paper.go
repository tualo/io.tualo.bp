package paper

import (
	"image"
	"math"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
)

type CC struct {
	Points gocv.PointVector
	Mat    cv.Mat
}
type Paper struct {
	contourChannel chan CC
	debugChannel   chan cv.Mat
}

func (me *Paper) Init() {
	me.contourChannel = make(chan CC, 1)
	me.debugChannel = make(chan cv.Mat, 1)
}

func (me *Paper) ContourChannel() chan CC {
	return me.contourChannel
}

func (me *Paper) DebugChannel() chan cv.Mat {
	return me.debugChannel
}

func (me *Paper) EdgeDifference(points map[string]image.Point) float64 {
	p1 := points["bottomLeftCorner"]
	p2 := points["bottomRightCorner"]
	p3 := points["topLeftCorner"]
	p4 := points["topRightCorner"]

	// horizontal
	d1 := math.Sqrt(float64((p1.X-p2.X)*(p1.X-p2.X)) + float64((p1.Y-p2.Y)*(p1.Y-p2.Y)))
	d2 := math.Sqrt(float64((p3.X-p4.X)*(p3.X-p4.X)) + float64((p3.Y-p4.Y)*(p3.Y-p4.Y)))

	// vertical
	d3 := math.Sqrt(float64((p4.X-p2.X)*(p4.X-p2.X)) + float64((p4.Y-p2.Y)*(p4.Y-p2.Y)))
	d4 := math.Sqrt(float64((p3.X-p1.X)*(p3.X-p1.X)) + float64((p3.Y-p1.Y)*(p3.Y-p1.Y)))

	f1 := 0.0
	f2 := 0.0
	if d1 > d2 {
		f1 = d1 / d2
	} else {
		f1 = d2 / d1
	}

	if d3 > d4 {
		f2 = d3 / d4
	} else {
		f2 = d4 / d3
	}

	return math.Max(f1, f2)
}
