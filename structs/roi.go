package structs

type ROI struct {
	X           float64
	Y           float64
	YCap        float64
	XCap        float64
	Height      float64
	Width       float64
	ItemCountX  int
	ItemCountY  int
	Ballotpaper int
	MaxAllowed  int
}

type ROIS struct {
	Roi []ROI
}
