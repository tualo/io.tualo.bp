package structs

type SendImageQueueItem struct {
	Barcode      string
	BoxBarcode   string
	StackBarcode string
	Id           int
	Marks        string
	Image        string
}
