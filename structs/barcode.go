package structs

import "image"

type BarcodeSymbol struct {
	Type     string
	Data     string
	Quality  int
	Boundary []image.Point
}
