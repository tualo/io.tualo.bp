package structs

import "gocv.io/x/gocv"

type SaveMarker struct {
	Mat        gocv.Mat
	Marked     string
	Pagination string
	R          int
	I          int
	P          int
	V          int
}
