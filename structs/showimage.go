package structs

import (
	"time"

	"tualo.de/deep-test/cv"
)

type ShowImageStruct struct {
	Title string
	Mat   cv.Mat
}

type LayerStruct struct {
	Title          string
	Mat            cv.Mat
	BasedOn        string
	DoNotDraw      bool
	Weighted_Alpha float64
	Weighted_Beta  float64
	Weighted_Gamma float64
	Timestamp      time.Time
}

type LayerRemoveStruct struct {
	Title     string
	Timestamp time.Time
}
