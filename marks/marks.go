package marks

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/globals"
	"tualo.de/deep-test/structs"
)

type Marks struct {
	mat                *gocv.Mat
	scale              structs.Scale
	InternalId         int
	bestContourTrigger func(gocv.Mat)
	showImg            func(gocv.Mat)
}

type MarkState struct {
	Found  bool
	Marked bool
	Value  float64
	Image  gocv.Mat
}

func (me *Marks) showImage(mat *gocv.Mat) {
	if me.showImg != nil {
		me.showImg(*mat)
	}
}

func (me *Marks) runBestContourTrigger(rect image.Rectangle, mat *gocv.Mat) {
	if me.bestContourTrigger != nil {
		m := mat.Region(rect)
		me.bestContourTrigger(m.Clone())
	}
}

func (me *Marks) SetBestContourTrigger(fn func(gocv.Mat)) {
	me.bestContourTrigger = fn
}

func (me *Marks) SetShowImageTrigger(fn func(gocv.Mat)) {
	me.showImg = fn
}

func (me *Marks) Set(mat *gocv.Mat, scale structs.Scale) {
	me.mat = mat
	me.scale = scale
}

/*
find the best contour, if possible with a circle contour
fit an rect over that contour and return that rectangle

important close the returned contour after use
*/
func (me *Marks) BestContour(frame *gocv.Mat, minArea float64) (image.Rectangle, gocv.PointsVector, int) {
	circleIndex := -1
	frameContour := 3
	// outer rect fix connected line on the image edges
	gocv.Rectangle(frame, image.Rectangle{image.Point{0, 0}, image.Point{frame.Cols(), frame.Rows()}}, color.RGBA{255, 255, 255, 120}, frameContour)
	frameX := frame.Clone()
	frame.ConvertTo(&frameX, gocv.MatTypeCV32SC1)
	cnts := gocv.FindContours(frameX, gocv.RetrievalFloodfill, gocv.ChainApproxNone)
	frameX.Close()
	var (
		bestCnt  gocv.PointVector
		bestArea = minArea
		found    = false
	)
	largest := 0.0
	for i := 0; i < cnts.Size(); i++ {
		if area := gocv.ContourArea(cnts.At(i)); area > largest {
			largest = (area)
		}
	}
	for i := 0; i < cnts.Size(); i++ {
		cnt := cnts.At(i)
		area := gocv.ContourArea(cnt)
		if area != largest {
			appoxCnt := gocv.ApproxPolyDP(cnt, 0.03*gocv.ArcLength(cnt, true), true)
			//log.Println("appoxCnt.Size()", appoxCnt.Size())
			if appoxCnt.Size() == 8 {
				_, _, radius := gocv.MinEnclosingCircle(appoxCnt)
				circleArea := radius * radius * math.Pi

				// log.Println("abweichung", math.Abs(math.Round(((float64(circleArea)/area)-1)*100)))
				if math.Abs(math.Round(((float64(circleArea)/area)-1)*100)) < 10 {
					circleIndex = i
				}
			}
			appoxCnt.Close()
			if area > bestArea {
				bestArea = area
				bestCnt = cnt
				found = true
			}
		}
	}
	// log.Println("circleIndex", circleIndex)
	if !found {
		return image.Rect(0, 0, frame.Cols(), frame.Rows()), cnts, circleIndex
	}

	pnts := bestCnt.ToPoints()
	x1 := 9999
	y1 := 9999
	x2 := 0
	y2 := 0
	for i := 0; i < len(pnts); i++ {
		x1 = int(math.Min(float64(x1), float64(pnts[i].X)))
		y1 = int(math.Min(float64(y1), float64(pnts[i].Y)))
		x2 = int(math.Max(float64(x2), float64(pnts[i].X)))
		y2 = int(math.Max(float64(y2), float64(pnts[i].Y)))
	}

	return image.Rect(x1, y1, x2, y2), cnts, circleIndex
}

func (me *Marks) Check() MarkState {

	thresholdMarks := 45.0 // je kleiner desto weniger rauschen
	dpHoughCircles := 1.0
	thresholdHoughCircles := 90.0
	accumulatorThresholdHoughCircles := 10.0

	settings := globals.Globals()

	if settings != nil {
		// thresholdMarks = settings.
		dpHoughCircles = settings.DpHoughCircles
		thresholdHoughCircles = settings.ThresholdHoughCircles
		accumulatorThresholdHoughCircles = settings.AccumulatorThresholdHoughCircles
	}

	circleSize := int(float64(4) * math.Min(me.scale.PixelScaleX, me.scale.PixelScaleY))
	minDist := (float64(12) * math.Min(me.scale.PixelScaleX, me.scale.PixelScaleY))
	strokeWidth := (float64(2) * math.Max(me.scale.PixelScaleX, me.scale.PixelScaleY))

	analyseMat := me.mat.Clone() // find best contour etc
	useResult := me.mat.Clone()  // find cross mark
	circles := gocv.NewMat()

	me.showImage(&useResult)
	defer analyseMat.Close()
	defer useResult.Close()

	gocv.HoughCirclesWithParams(
		analyseMat,
		&circles,
		gocv.HoughGradient,
		dpHoughCircles,                   // dp
		minDist,                          //float64(analyseMat.Rows()/50), // minDist
		thresholdHoughCircles,            // param1
		accumulatorThresholdHoughCircles, // param2
		circleSize,                       // minRadius
		circleSize,                       // maxRadius
	)
	defer circles.Close()

	// log.Println("circles ", circles.Cols(), "x", circles.Rows(), circleSize, thresholdMarks, globals.Globals())

	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])

			sx := x - r - int(me.scale.PixelScaleX*0.5)
			sy := y - r - int(me.scale.PixelScaleY*0.5)
			ex := x + r + int(me.scale.PixelScaleX*0.5)
			ey := y + r + int(me.scale.PixelScaleY*0.5)

			if sx < 0 {
				sx = 0
			}
			if sx > analyseMat.Cols() {
				// ToDo TriState: No Circle Found
				return MarkState{
					Marked: false,
					Found:  false,
					Value:  -1,
				}
			}
			if sy < 0 {
				sy = 0
			}
			if sy > analyseMat.Rows() {
				// ToDo TriState: No Circle Found
				return MarkState{
					Marked: false,
					Found:  false,
					Value:  -1,
				}
			}

			if ex > analyseMat.Cols() {
				ex = analyseMat.Cols()
			}
			if ey > analyseMat.Rows() {
				ey = analyseMat.Rows()
			}

			rect := image.Rect(
				sx,
				sy,
				ex,
				ey,
			)

			gocv.Threshold(analyseMat, &analyseMat, float32(thresholdMarks), 255, gocv.ThresholdBinary)

			// remove noise
			gocv.Erode(analyseMat, &analyseMat, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(7, 7)))
			gocv.Dilate(analyseMat, &analyseMat, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(5, 5)))

			analyseMat = analyseMat.Region(rect)
			me.showImage(&analyseMat)
			useResult = useResult.Region(rect)
			best, contours, circleIndex := me.BestContour(&analyseMat, float64(r*r)*0.5)

			me.runBestContourTrigger(best, &useResult)

			// gocv.IMWrite(fmt.Sprintf("marks_line_runBestContourTrigger_%d.jpg", me.InternalId), useResult)

			// wenn keine kreiskontur gefunden wurde (-1),
			// gibt es verbindungen zum kreis > ergo eine markierung
			if circleIndex > -1 {
				gocv.DrawContours(&useResult, contours, circleIndex, color.RGBA{255, 255, 255, 0}, int(strokeWidth))
				// gocv.IMWrite(fmt.Sprintf("marks_line_no_circle_%d.jpg", me.InternalId), useResult)

			} else {
				replacementCircles := gocv.NewMat()
				gocv.HoughCirclesWithParams(
					useResult,
					&replacementCircles,
					gocv.HoughGradient,
					dpHoughCircles,                   // dp
					minDist,                          //float64(analyseMat.Rows()/50), // minDist
					thresholdHoughCircles,            // param1
					accumulatorThresholdHoughCircles, // param2
					circleSize,                       // minRadius
					circleSize,                       // maxRadius
				)

				for i := 0; i < replacementCircles.Cols(); i++ {
					v := replacementCircles.GetVecfAt(0, i)
					if len(v) > 2 {
						x := int(v[0])
						y := int(v[1])
						r := int(v[2])
						gocv.Circle(&useResult, image.Pt(x, y), r-3, color.RGBA{255, 255, 255, 0}, int(strokeWidth))
						gocv.Circle(&useResult, image.Pt(x, y), r-2, color.RGBA{255, 255, 255, 0}, int(strokeWidth))
						gocv.Circle(&useResult, image.Pt(x, y), r-1, color.RGBA{255, 255, 255, 0}, int(strokeWidth))
						gocv.Circle(&useResult, image.Pt(x, y), r-0, color.RGBA{255, 255, 255, 0}, int(strokeWidth))
						gocv.Circle(&useResult, image.Pt(x, y), r+int(strokeWidth), color.RGBA{255, 255, 255, 0}, int(strokeWidth)*3)
						gocv.Circle(&useResult, image.Pt(x, y), r+int(strokeWidth)+int(strokeWidth), color.RGBA{255, 255, 255, 0}, int(strokeWidth)*3)
						gocv.Circle(&useResult, image.Pt(x, y), r+int(strokeWidth)+int(strokeWidth)+int(strokeWidth), color.RGBA{255, 255, 255, 0}, int(strokeWidth)*3)
						gocv.Circle(&useResult, image.Pt(x, y), r+int(strokeWidth)+int(strokeWidth)+int(strokeWidth), color.RGBA{255, 255, 255, 0}, int(strokeWidth)*3)
					}
				}
				replacementCircles.Close()
				// gocv.IMWrite(fmt.Sprintf("marks_line_replacementCircles_%d.jpg", me.InternalId), useResult)

			}
			gocv.Erode(useResult, &useResult, gocv.GetStructuringElement(gocv.MorphEllipse, image.Pt(3, 3)))

			// gocv.IMWrite(fmt.Sprintf("marks_line_erode_%d.jpg", me.InternalId), useResult)

			analyseMat = analyseMat.Region(best)
			useResult = useResult.Region(best)

			gocv.Threshold(useResult, &useResult, float32(thresholdMarks), 255, gocv.ThresholdBinary)
			gocv.BitwiseNot(useResult, &useResult)
			me.showImage(&useResult)
			cx := (float64(useResult.Cols()) * float64(useResult.Rows()))
			avg := ((useResult.Sum().Val1) / cx) / 255.0

			/*
				if (useResult.Sum().Val1 / cx) > 1 {
					gocv.IMWrite(fmt.Sprintf("marks_line_useresult_%d.jpg", me.InternalId), useResult)
				}
			*/

			// noch nicht sicher,
			// ob das so richtig ist, aber es funktioniert
			return MarkState{
				Marked: (useResult.Sum().Val1 / cx) > 1,
				Found:  true,
				Value:  avg,
				Image:  useResult,
			}
		}
	}
	return MarkState{
		Marked: false,
		Found:  false,
		Value:  -1,
	}

}
