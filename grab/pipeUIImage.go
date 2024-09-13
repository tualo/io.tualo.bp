package grab

import (
	"image"
	"gocv.io/x/gocv"
)

func (this *GrabcameraClass) pipeUIImage(img gocv.Mat){
	if len(this.imageChannelPaper)==cap(this.imageChannelPaper) {
		mat,_:=<-this.imageChannelPaper
		mat.Close()
	}

	cloned := img.Clone()
	gocv.Resize(cloned, &cloned, image.Point{}, 0.3, 0.3 , gocv.InterpolationLinear)
	this.imageChannelPaper <- cloned
	//cloned.Close()
}