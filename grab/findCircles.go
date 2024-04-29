package grab

import (
	"image"
 	"sort"
	"math"
	"log"
	"image/color"
	"gocv.io/x/gocv"
	structs "io.tualo.bp/structs"
)


func DrawCircles(img *gocv.Mat, circles *gocv.Mat,  innerOverdraw int, outerOverdraw int, marks []structs.CheckMarks ) {
	var _color color.RGBA = color.RGBA{255, 255, 255, 0}
	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])
			if r-innerOverdraw/10> 0 {
				if len(marks) > i {
					_color = color.RGBA{220, 220, 220, 0}
					/*
					if math.Round(marks[i].Mean) > meanFindCircles  {
						_color = color.RGBA{0, 255, 0, 0}
					}else{
						_color = color.RGBA{220, 220, 220, 0}
					}
					*/
				}
				gocv.Circle(img, image.Pt(x, y), r-innerOverdraw/10, _color, outerOverdraw/10)
			}
		}
	}
}

func (this *GrabcameraClass) findCircles(croppedMat gocv.Mat , circleSize int,minDist float64) []structs.CheckMarks {
	croppedMatGray := gocv.NewMat()
	gocv.CvtColor(croppedMat, &croppedMatGray, gocv.ColorBGRToGray)
	circles := gocv.NewMat()
	if this.globals.DpHoughCircles == 0 {
		this.globals.SetDefaults()
	}
	
	gocv.HoughCirclesWithParams(
		croppedMatGray,
		&circles,
		gocv.HoughGradient,
		this.globals.DpHoughCircles,                     // dp
		minDist, //float64(croppedMatGray.Rows()/50), // minDist
		this.globals.ThresholdHoughCircles,                    // param1
		this.globals.AccumulatorThresholdHoughCircles,                    // param2
		circleSize,                    // minRadius
		circleSize,                     // maxRadius
	)

	this.globals.InnerOverdrawDrawCircles = 2

	log.Println("circles",circleSize, this.globals.InnerOverdrawDrawCircles*int(this.pixelScale), this.globals.GaussianBlurFindCircles/ int(this.pixelScale) )

	imgRGray := gocv.NewMat()
	imgGray := gocv.NewMat()
	imgBlur := gocv.NewMat()
	gocv.CvtColor(croppedMat, &imgGray, gocv.ColorBGRToGray)


	blurSize := int(math.Round( float64(this.globals.GaussianBlurFindCircles) * this.pixelScale))
	if blurSize % 2 == 0 {
		blurSize++
	}

	gocv.GaussianBlur(imgGray, &imgBlur, image.Point{blurSize, blurSize}, 0, 0, gocv.BorderDefault)
	gocv.AdaptiveThreshold(imgBlur, &imgRGray, 255.0, gocv.AdaptiveThresholdGaussian, gocv.ThresholdBinary, this.globals.AdaptiveThresholdBlockSize, this.globals.AdaptiveThresholdSubtractMean)
	imgBlur.Close()
	imgGray.Close()

	checkMarks := []structs.CheckMarks{}
	//checkMarksList := []bool{}


	DrawCircles(&imgRGray, &circles,  this.globals.InnerOverdrawDrawCircles*int(this.pixelScale), this.globals.OuterOverdrawDrawCircles*int(this.pixelScale), checkMarks)
	this.globals.MeanFindCircles = 253
	for i := 0; i < circles.Cols(); i++ {
		v := circles.GetVecfAt(0, i)
		if len(v) > 2 {
			x := int(v[0])
			y := int(v[1])
			r := int(v[2])
			rect_circle:=image.Rect(x-r , y-r  , x+r , y+r )
			if rect_circle.Min.X < 0 || rect_circle.Min.Y < 0 || rect_circle.Max.X > imgRGray.Cols() || rect_circle.Max.Y > imgRGray.Rows() {
				continue
			}else{
				rect_circleMat := imgRGray.Region(rect_circle)
				mean := rect_circleMat.Mean()
				rect_circleMat.Close()
				checkMarks = append(checkMarks, structs.CheckMarks{mean.Val1, x, y, r,math.Round( mean.Val1 ) < this.globals.MeanFindCircles})
				log.Println("check  ",i,this.globals.MeanFindCircles,math.Round( mean.Val1 ) < this.globals.MeanFindCircles)
			}
		}
	}

	log.Println("checkMarks",checkMarks,this.globals.InnerOverdrawDrawCircles*int(this.pixelScale), this.globals.OuterOverdrawDrawCircles*int(this.pixelScale))

	sort.Slice(checkMarks[:], func(i, j int) bool {
		return checkMarks[i].Y < checkMarks[j].Y
	})
	imgCol := gocv.NewMat()
	gocv.CvtColor(imgRGray, &imgCol, gocv.ColorGrayToBGR)
	DrawCircles(&imgCol, &circles, this.globals.InnerOverdrawDrawCircles*int(this.pixelScale), this.globals.OuterOverdrawDrawCircles*int(this.pixelScale), checkMarks)

	//gocv.IMWrite("circles.jpg", imgRGray)

	circles.Close()
	imgRGray.Close()

	return checkMarks
}