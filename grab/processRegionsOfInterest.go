package grab

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) processRegionsOfInterest701(tr structs.TesseractReturnType, img gocv.Mat, useRois []int) {
	pImage := img.Clone()

	for useRoi := 0; useRoi < len(useRois); useRoi++ {
		if useRoi < len(tr.PageRois) {
			// for pRoiIndex := 0; pRoiIndex < len(tr.PageRois); pRoiIndex++ {
			X := int(float64(tr.PageRois[useRois[useRoi]].X) * this.pixelScale)
			Y := int(float64(tr.PageRois[useRois[useRoi]].Y) * this.pixelScaleY)
			W := int(float64(tr.PageRois[useRois[useRoi]].Width) * this.pixelScale)
			H := int(float64(tr.PageRois[useRois[useRoi]].Height) * this.pixelScaleY)
			rect := image.Rect(X, Y, X+W, Y+H)
			gocv.Rectangle(&pImage, rect, color.RGBA{0, 0, 255, 125}, int(10*this.pixelScaleY))
			//DrawCircles(&imgCol, &circles, this.globals.InnerOverdrawDrawCircles*int(this.pixelScale), this.globals.OuterOverdrawDrawCircles*int(this.pixelScale), checkMarks)
		}
	}
	this.pipeUIImage(pImage)
	pImage.Close()
}

func (this *GrabcameraClass) processRegionsOfInterest(tr structs.TesseractReturnType, img gocv.Mat, useRois []int) structs.TesseractReturnType {

	this.pixelScale = float64(img.Cols()) / float64(tr.Pagesize.Width)
	this.pixelScaleY = float64(img.Rows()) / float64(tr.Pagesize.Height)

	if this.pixelScale == 0 {
		this.pixelScale = 1
	}
	if this.pixelScaleY == 0 {
		this.pixelScaleY = 1
	}
	circleSize := int(float64(tr.CircleSize) * this.pixelScale)
	minDist := float64(tr.CircleMinDistance) * this.pixelScale

	marks := []structs.CheckMarks{}
	for useRoi := 0; useRoi < len(useRois); useRoi++ {
		if useRoi < len(tr.PageRois) {
			// for pRoiIndex := 0; pRoiIndex < len(tr.PageRois); pRoiIndex++ {
			X := int(float64(tr.PageRois[useRois[useRoi]].X) * this.pixelScale)
			Y := int(float64(tr.PageRois[useRois[useRoi]].Y) * this.pixelScaleY)
			W := int(float64(tr.PageRois[useRois[useRoi]].Width) * this.pixelScale)
			H := int(float64(tr.PageRois[useRois[useRoi]].Height) * this.pixelScaleY)

			rect := image.Rect(X, Y, X+W, Y+H)
			croppedMat := img.Region(rect)

			if !croppedMat.Empty() {
				fMarks := this.findCircles(croppedMat, circleSize, minDist, useRois[useRoi])
				for i := 0; i < len(fMarks); i++ {
					marks = append(marks, fMarks[i])
				}
			}
			croppedMat.Close()
		}
	}

	if this.globals.ShowImage == 701 {
		this.processRegionsOfInterest701(tr, img, useRois)
	}

	tr.Marks = marks
	if tr.PageRois[useRois[0]].ExcpectedMarks == len(marks) {
		tr.IsCorrect = true
	}

	return tr

}
