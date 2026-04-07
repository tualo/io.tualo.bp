package structs

import "image/color"

type HistoryListItem struct {
	Barcode      string
	BoxBarcode   string
	StackBarcode string
	State        string
	StateColor   color.RGBA
}
