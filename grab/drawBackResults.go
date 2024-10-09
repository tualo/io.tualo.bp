package grab

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
	// "log"
)

func (this *GrabcameraClass) drawBackResults() {

	/*
		for i := 0; i < len(this.debugMarkList); i++ {
			//playGround
			if this.debugMarkList[i].Checked {
				gocv.Circle(&this.playGround, this.debugMarkList[i].Point, this.debugMarkList[i].Pixelsize, color.RGBA{100, 100, 100, 120}, int(1.0*this.pixelScale))
			} else {
				gocv.Circle(&this.playGround, this.debugMarkList[i].Point, this.debugMarkList[i].Pixelsize, color.RGBA{100, 100, 100, 120}, int(2.0*this.pixelScale))
			}
		}
	*/

	// Drawing back the results
	for i := 0; i < len(this.checkMarkList); i++ {
		//playGround
		if this.checkMarkList[i].Checked {
			gocv.Circle(&this.playGround, this.checkMarkList[i].Point, this.checkMarkList[i].Pixelsize, color.RGBA{0, 255, 0, 120}, int(3.0*this.pixelScale))
		} else {
			gocv.Circle(&this.playGround, this.checkMarkList[i].Point, this.checkMarkList[i].Pixelsize, color.RGBA{255, 0, 0, 120}, int(3.0*this.pixelScale))
		}
	}

	gocv.WarpPerspective(this.playGround, &this.img, this.invM, image.Point{this.img.Cols(), this.img.Rows()})
	if this.contour.Size() != 0 {
		drawContours := gocv.NewPointsVector()
		drawContours.Append(this.contour)
		// log.Println("drawBackResults drawContours",uint8(this.currentState.Red), uint8(this.currentState.Green), uint8(this.currentState.Blue))
		gocv.DrawContours(&this.img, drawContours, -1, color.RGBA{uint8(this.currentState.Red), uint8(this.currentState.Green), uint8(this.currentState.Blue), 120}, int(8.0*this.pixelScale))
		drawContours.Close()
	}
}
