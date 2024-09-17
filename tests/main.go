
package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"sort"

	"gocv.io/x/gocv"
)

type CustomContour struct {
	c gocv.PointsVector
	index int
	area float64
}

func (cc CustomContour) Len() int {
	return cc.c.Size()
}

func (cc CustomContour) Less(i, j int) bool {
	aI := gocv.ContourArea(cc.c.At(i))
	aJ := gocv.ContourArea(cc.c.At(j))
	if aI > aJ {
		return true
	}
	return false
}

func (cc CustomContour) Swap(i, j int) {
	/*
	oI=cc.c.At(i)
	oJ=cc.c.At(j)
	cc.c.Set(i, oJ)
	cc.c.Set(j, oI)
	//cc.c.At(i), cc.c.At(j) = cc.c.At(j), cc.c.At(i)
	*/
	points := cc.c.ToPoints()
	points[i], points[j] =points[j],points[i]
	cc.c = gocv.NewPointsVectorFromPoints(points)
}


func main() {
	filename := os.Args[1]

	mat := gocv.IMRead(filename, gocv.IMReadColor)

	matCanny := gocv.NewMat()
	matLines := gocv.NewMat()

	window := gocv.NewWindow("detected lines")

	kernel :=  gocv.Ones(5,5,gocv.MatTypeCV8U) // np.ones((5,5),np.uint8)
	gocv.MorphologyExWithParams(mat, &mat, gocv.MorphErode, kernel, 3,gocv.BorderDefault)
	morph_window := gocv.NewWindow("morph_window")
	morph_window.IMShow(mat)

	gocv.Canny(mat, &matCanny, 50, 250)
	canny_window := gocv.NewWindow("canny_window")
	canny_window.IMShow(matCanny)



	//add some gausian blue to make edge detection easier
	blur := gocv.NewMat()
	gocv.GaussianBlur(mat, &blur, image.Pt(3, 3), 1, 1, gocv.BorderDefault)

	//erode edges so that 'more is inside of the edge detection
	eroded := gocv.NewMat()
	{
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(20, 20))
		defer kernel.Close()
		gocv.Erode(blur, &eroded, kernel)
	}
	//median blur to remove 'salt and pepper'
	medianBlur := gocv.NewMat()
	gocv.MedianBlur(eroded, &medianBlur, 9)
	eroded_window := gocv.NewWindow("eroded_window")
	eroded_window.IMShow(medianBlur)


	//use morphology to try and join lines up to create a continuous line
	morph := gocv.NewMat()
	{
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()

		gocv.MorphologyEx(medianBlur, &morph, gocv.MorphClose, kernel)
	}

	contours := gocv.FindContours(matCanny, gocv.RetrievalExternal, gocv.ChainApproxSimple)

	var toSort CustomContour
	toSort.c = contours
	//find the contour with the largest area
	sort.Sort(CustomContour(toSort))

	statusColor := color.RGBA{255, 0, 0, 0}

	if toSort.c.Size() > 0 {
		gocv.FillPoly(&mat, toSort.c, statusColor)
		polyfill_window := gocv.NewWindow("polyfill_window")
		polyfill_window.IMShow(mat)
		}


	/*
	mask = np.zeros(img.shape[:2],np.uint8)
bgdModel = np.zeros((1,65),np.float64)
fgdModel = np.zeros((1,65),np.float64)
rect = (20,20,img.shape[1]-20,img.shape[0]-20)
	cv2.grabCut(img,mask,rect,bgdModel,fgdModel,5,cv2.GC_INIT_WITH_RECT)
	*/
	/*
	mask := gocv.NewMatWithSize(mat.Rows(), mat.Cols(), gocv.MatTypeCV8U)
	bgdModel := gocv.Zeros(1, 65, gocv.MatTypeCV64F)
	fgdModel := gocv.Zeros(1, 65, gocv.MatTypeCV64F)
	rect := image.Rect(20, 20, mat.Cols()-20, mat.Rows()-20)
	gocv.GrabCut(mat, &mask, rect, &bgdModel, &fgdModel, 5, gocv.GCInitWithMask)
	*/
	//gocv.GrabCut(img Mat, mask *Mat, r image.Rectangle, bgdModel *Mat, fgdModel *Mat, iterCount int, mode GrabCutMode)



	gocv.HoughLinesP(matCanny, &matLines, 1, math.Pi/180, 80)

	fmt.Println(matLines.Cols())
	fmt.Println(matLines.Rows())
	for i := 0; i < matLines.Rows(); i++ {
		pt1 := image.Pt(int(matLines.GetVeciAt(i, 0)[0]), int(matLines.GetVeciAt(i, 0)[1]))
		pt2 := image.Pt(int(matLines.GetVeciAt(i, 0)[2]), int(matLines.GetVeciAt(i, 0)[3]))
		gocv.Line(&mat, pt1, pt2, color.RGBA{0, 255, 0, 50}, 10)
	}

	for {
		window.IMShow(mat)
		if window.WaitKey(10) >= 0 {
			break
		}
	}
}
