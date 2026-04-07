package background

import (
	"gocv.io/x/gocv"
	"tualo.de/deep-test/cv"
)

type Background struct {
	mog           gocv.BackgroundSubtractorKNN
	imgDelta      cv.Mat
	imgThresh     cv.Mat
	resultChannel chan cv.Mat
}

func (me *Background) ResultChannel() chan cv.Mat {
	return me.resultChannel
}

func (me *Background) Init() {
	me.mog = gocv.NewBackgroundSubtractorKNNWithParams(20, 240, true)
	me.imgDelta = cv.NewMat("Background imgDelta")
	me.imgThresh = cv.NewMat("Background imgThresh")
	me.resultChannel = make(chan cv.Mat, 1)
	// defer me.mog.Close()
}

func (me *Background) Check(img cv.Mat) {
	me.mog.Apply(img.Get(), me.imgDelta.GetPointer())
	if len(me.resultChannel) == cap(me.resultChannel) {
		mat := <-me.resultChannel
		mat.Close()
	}
	me.resultChannel <- me.imgDelta.Clone("Background mask")
}
