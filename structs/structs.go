package structs

import (
	"image"
	"image/color"
	"gocv.io/x/gocv"
)

type TesseractReturnType struct {
	Point    image.Point
	BoxBarcode   string
	StackBarcode   string
	Barcode   string
	Title   string
	Id string
	IsCorrect bool
	Marks   []CheckMarks
	PageRois []DocumentConfigurationPageRoi
	Pagesize DocumentConfigurationPageSize
	CircleSize int
	CircleMinDistance int
}

type RoisChannelStruct struct {
	tesseractReturn TesseractReturnType
	mat gocv.Mat
}

type CheckMarkList struct {
	Count int
	Sum int
	AVG float64
	Checked bool
	Point image.Point
	Pixelsize int
}

type CheckMarks struct {
	Mean float64
    X       int 
	Y       int
	Radius	   int
	Checked bool
	RoiIndex int
}


type ReturnType struct {
	Point    image.Point
	FCBarcode   string
	Barcode   string
	Title   string
	IsCorrect bool
	Marks   []bool
	PageRois []DocumentConfigurationPageRoi
	Pagesize DocumentConfigurationPageSize
	CircleSize int
	CircleMinDistance int
}

type CameraList struct {
	Width int
	Height int
	Index int
	Title string
}



type DocumentConfigurationPageSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type DocumentConfigurationPageRoi struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
	ExcpectedMarks int `json:"excpectedMarks"`
	Types	[]struct{
		Title string `json:"title"`
		Id int `json:"id"`
	} `json:"types"`
	
}

type DocumentConfigurations []struct {
	Titles       []string `json:"titles"`
	CircleSize   int `json:"circleSize"`
	CircleMinDistance int `json:"circleMinDistance"`
	TitleRegion struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"titleRegion"`
	Pagesize DocumentConfigurationPageSize `json:"pagesize"`
	Rois []DocumentConfigurationPageRoi `json:"rois"`
}

type BarcodeSymbol struct {
	Type string
	Data string
	Quality int
	Boundary []image.Point
}






type ImageProcessorState struct {
	Name string
	Red int
	Green int
	Blue int
	Opacity int
}

type HistoryListItem struct {
	Barcode string
	BoxBarcode string
	StackBarcode string
	State string
	StateColor color.RGBA
}


type DetectedCodes struct {
	Barcode string
	BoxBarcode string
	StackBarcode string
}