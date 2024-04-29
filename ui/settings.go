package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fmt"
	globals "io.tualo.bp/globals"
	grabcamera "io.tualo.bp/grab"
)

type SettingsScreenClass struct {
	grabber *grabcamera.GrabcameraClass
	globals *globals.GlobalValuesClass
	cameraSelectWidget *widget.Select
	thresholdHoughCirclesWidget *widget.Slider
	meanFindCirclesWidget *widget.Slider
	dpHoughCirclesWidget *widget.Slider
	gaussianBlurFindCirclesWidget *widget.Slider
	adaptiveThresholdBlockSizeWidget *widget.Slider
	adaptiveThresholdSubtractMeanWidget *widget.Slider
	thresholdHoughCirclesWidgetLabel *widget.Label
	meanFindCirclesWidgetLabel *widget.Label
	dpHoughCirclesWidgetLabel *widget.Label
	gaussianBlurFindCirclesWidgetLabel *widget.Label
	adaptiveThresholdBlockSizeWidgetLabel *widget.Label
	adaptiveThresholdSubtractMeanWidgetLabel *widget.Label
	paperImageCheck *widget.Check
}
func (this *SettingsScreenClass) SetGrabber(grabber *grabcamera.GrabcameraClass) {
	this.grabber = grabber
}
func (this *SettingsScreenClass) SetGlobals(globals *globals.GlobalValuesClass) {
	this.globals = globals
}
func (this *SettingsScreenClass) makeSettingsForm() fyne.CanvasObject {

	cameraList := this.grabber.GetCameraList()
	fmt.Println("maxcameranum",len(cameraList))	
	//"Camera 1", "Camera 2", "Camera 3", "Camera 4"
	this.cameraSelectWidget = widget.NewSelect([]string{},func(value string) {
		fmt.Println("cameraSelectWidget",value)
		for i:=0;i<len(cameraList);i++ {
			if value == fmt.Sprintf("Camera %d (%dx%d)",(i+1),cameraList[i].Width,cameraList[i].Height) {
				this.globals.IntCamera = i
			}
		}
	})
	for i:=0;i<len(cameraList);i++ {
		this.cameraSelectWidget.Options = append(this.cameraSelectWidget.Options,fmt.Sprintf("Camera %d (%dx%d)",(i+1),cameraList[i].Width,cameraList[i].Height))
	}
	this.cameraSelectWidget.PlaceHolder = "Bitte wählen Sie eine Kamera aus"
	fmt.Println("cameraSelectWidget",this.cameraSelectWidget.Options,cameraList,this.globals.IntCamera)
	if this.globals.IntCamera>len(cameraList) {
		this.globals.IntCamera = 0
	}
	if len(cameraList)>0 {
		this.cameraSelectWidget.SetSelected(this.cameraSelectWidget.Options[this.globals.IntCamera])
	}


	this.thresholdHoughCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", this.globals.ThresholdHoughCircles))
	this.meanFindCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", this.globals.MeanFindCircles))
	this.dpHoughCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", this.globals.DpHoughCircles))
	this.gaussianBlurFindCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%d", this.globals.GaussianBlurFindCircles))
	this.adaptiveThresholdBlockSizeWidgetLabel = widget.NewLabel(fmt.Sprintf("%d", this.globals.AdaptiveThresholdBlockSize))
	this.adaptiveThresholdSubtractMeanWidgetLabel = widget.NewLabel(fmt.Sprintf("%.1f", this.globals.AdaptiveThresholdSubtractMean))



	this.thresholdHoughCirclesWidget = widget.NewSlider(0, 255)
	this.thresholdHoughCirclesWidget.Value = this.globals.ThresholdHoughCircles
	this.thresholdHoughCirclesWidget.OnChangeEnded = func(value float64) {
		this.globals.ThresholdHoughCircles = value
		this.thresholdHoughCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}


	this.meanFindCirclesWidget = widget.NewSlider(1, 254)
	this.meanFindCirclesWidget.Value = this.globals.MeanFindCircles
	this.meanFindCirclesWidget.OnChangeEnded = func(value float64) {
		this.globals.MeanFindCircles = value
		this.meanFindCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}

	this.dpHoughCirclesWidget = widget.NewSlider(0, 3)
	this.dpHoughCirclesWidget.Value = this.globals.DpHoughCircles
	this.dpHoughCirclesWidget.OnChangeEnded = func(value float64) {
		this.globals.DpHoughCircles = value
		this.dpHoughCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}

	this.gaussianBlurFindCirclesWidget = widget.NewSlider(0, 255)
	this.gaussianBlurFindCirclesWidget.Value = float64(this.globals.GaussianBlurFindCircles)
	this.gaussianBlurFindCirclesWidget.OnChangeEnded = func(value float64) {

		this.globals.GaussianBlurFindCircles = int(value)
		this.gaussianBlurFindCirclesWidgetLabel.SetText(fmt.Sprintf("%d", int(value)))

	}

	this.adaptiveThresholdBlockSizeWidget = widget.NewSlider(0, 255)
	this.adaptiveThresholdBlockSizeWidget.Value = float64(this.globals.AdaptiveThresholdBlockSize)
	this.adaptiveThresholdBlockSizeWidget.OnChangeEnded = func(value float64) {
		this.globals.AdaptiveThresholdBlockSize = int(value)
		this.adaptiveThresholdBlockSizeWidgetLabel.SetText(fmt.Sprintf("%d", int(value)))

	}

	this.adaptiveThresholdSubtractMeanWidget = widget.NewSlider(-10, 10)
	this.adaptiveThresholdSubtractMeanWidget.Value = float64(this.globals.AdaptiveThresholdSubtractMean)
	this.adaptiveThresholdSubtractMeanWidget.OnChangeEnded = func(value float64) {
		this.globals.AdaptiveThresholdSubtractMean = float32(value)
		this.adaptiveThresholdSubtractMeanWidgetLabel.SetText(fmt.Sprintf("%.1f", value))
	}
	txt:=widget.NewLabel("Camera")
	
	container := container.New(layout.NewVBoxLayout(), 
	txt,
	this.cameraSelectWidget,
	widget.NewAccordion(
		&widget.AccordionItem{
			Title:  "Kreisdetetion",
			Detail: container.New(
				layout.NewGridLayout(1), 
				widget.NewLabel("Mean Find Circles"),
				container.NewBorder( nil, nil, nil, this.meanFindCirclesWidgetLabel,this.meanFindCirclesWidget ),

				widget.NewLabel("Hough Circles Threshold"),
				container.NewBorder( nil, nil, nil, this.thresholdHoughCirclesWidgetLabel,this.thresholdHoughCirclesWidget ),

				widget.NewLabel("Inverse ratio of the accumulator"),
				container.NewBorder( nil, nil, nil, this.dpHoughCirclesWidgetLabel,this.dpHoughCirclesWidget ),

				widget.NewLabel("Blursize"),
				container.NewBorder( nil, nil, nil, this.gaussianBlurFindCirclesWidgetLabel,this.gaussianBlurFindCirclesWidget ),


				widget.NewLabel("Adaptive Threshold Block Size"),
				container.NewBorder( nil, nil, nil, this.adaptiveThresholdBlockSizeWidgetLabel,this.adaptiveThresholdBlockSizeWidget ),

				widget.NewLabel("Adaptive Threshold Subtract Mean"),
				container.NewBorder( nil, nil, nil, this.adaptiveThresholdSubtractMeanWidgetLabel,this.adaptiveThresholdSubtractMeanWidget ), 
			), 
		}, 
	),
)
	return  container
}


func (t *SettingsScreenClass) CreateContainer() *fyne.Container {
	container := container.New(
		layout.NewPaddedLayout(),
		t.makeSettingsForm(),
	)
	return container
}

func NewSettingsScreenClass() *SettingsScreenClass {
	o := &SettingsScreenClass{
	}
	return o
}