package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"tualo.de/deep-test/camera"
	"tualo.de/deep-test/globals"
)

type SettingsScreenClass struct {
	/*
		grabber *grabcamera.GrabcameraClass
		globals *globals.GlobalValuesClass
	*/

	cameraSelectWidget                  *widget.Select
	cameraCaptureFrameFactorWidget      *widget.Slider
	cameraCaptureFrameFactorWidgetLabel *widget.Label

	cameraCaptureFPSWidget      *widget.Slider
	cameraCaptureFPSWidgetLabel *widget.Label

	paperContourFactorWidget      *widget.Slider
	paperContourFactorWidgetLabel *widget.Label

	thresholdHoughCirclesWidget              *widget.Slider
	meanFindCirclesWidget                    *widget.Slider
	dpHoughCirclesWidget                     *widget.Slider
	gaussianBlurFindCirclesWidget            *widget.Slider
	adaptiveThresholdBlockSizeWidget         *widget.Slider
	adaptiveThresholdSubtractMeanWidget      *widget.Slider
	thresholdHoughCirclesWidgetLabel         *widget.Label
	meanFindCirclesWidgetLabel               *widget.Label
	dpHoughCirclesWidgetLabel                *widget.Label
	gaussianBlurFindCirclesWidgetLabel       *widget.Label
	adaptiveThresholdBlockSizeWidgetLabel    *widget.Label
	adaptiveThresholdSubtractMeanWidgetLabel *widget.Label
	paperImageCheck                          *widget.Check
	checkR                                   *widget.Check
	checkG                                   *widget.Check
	checkB                                   *widget.Check
	paperFindContourNoiseBlurSizeWidget      *widget.Slider
	paperFindContourNoiseBlurSizeWidgetLabel *widget.Label
	erodeDillateSizeWidget                   *widget.Slider
	erodeDillateSizeWidgetLabel              *widget.Label

	outputWidget *widget.Select
}

func (this *SettingsScreenClass) SetFCChannels(value bool) {
	val := 0
	if this.checkR.Checked {
		val |= 1 << 2
	}
	if this.checkG.Checked {
		val |= 1 << 1
	}
	if this.checkB.Checked {
		val |= 1 << 0
	}
	globals.Globals().FindContourChannelMask = val
}

func (this *SettingsScreenClass) makeSettingsForm() fyne.CanvasObject {

	cameraList := camera.GetCameraList()
	// fmt.Println("maxcameranum",len(cameraList))
	//"Camera 1", "Camera 2", "Camera 3", "Camera 4"
	this.cameraSelectWidget = widget.NewSelect([]string{}, func(value string) {
		// fmt.Println("cameraSelectWidget",value)
		for i := 0; i < len(cameraList); i++ {
			if value == fmt.Sprintf("Camera %d (%dx%d)", (i+1), cameraList[i].Width, cameraList[i].Height) {
				globals.Globals().IntCamera = i
			}
		}
	})
	for i := 0; i < len(cameraList); i++ {
		this.cameraSelectWidget.Options = append(this.cameraSelectWidget.Options, fmt.Sprintf("Camera %d (%dx%d)", (i+1), cameraList[i].Width, cameraList[i].Height))
	}
	this.cameraSelectWidget.PlaceHolder = "Bitte wählen Sie eine Kamera aus"
	// fmt.Println("cameraSelectWidget",this.cameraSelectWidget.Options,cameraList,globals.Globals().IntCamera)
	if globals.Globals().IntCamera > len(cameraList) {
		globals.Globals().IntCamera = 0
	}
	if len(cameraList) > 0 {
		if globals.Globals().IntCamera < len(cameraList) {
			this.cameraSelectWidget.SetSelected(this.cameraSelectWidget.Options[globals.Globals().IntCamera])
		} else {
			this.cameraSelectWidget.SetSelected(this.cameraSelectWidget.Options[0])
		}
	}

	this.cameraCaptureFrameFactorWidgetLabel = widget.NewLabel(fmt.Sprintf("%.2f", globals.Globals().CaptureFrameFactor))
	this.cameraCaptureFrameFactorWidget = widget.NewSlider(1, 3)
	// 1 = 1.0, 2 = 0.75, 3 = 0.5
	if globals.Globals().CaptureFrameFactor == 1 {
		this.cameraCaptureFrameFactorWidget.Value = 1
	} else if globals.Globals().CaptureFrameFactor == 0.75 {
		this.cameraCaptureFrameFactorWidget.Value = 2
	} else if globals.Globals().CaptureFrameFactor == 0.5 {
		this.cameraCaptureFrameFactorWidget.Value = 3
	}
	// this.cameraCaptureFrameFactorWidget.Value = 1 / globals.Globals().CaptureFrameFactor
	this.cameraCaptureFrameFactorWidget.OnChangeEnded = func(value float64) {
		fmt.Println("cameraCaptureFrameFactorWidget", value)
		if value == 1 {
			globals.Globals().CaptureFrameFactor = 1
		} else if value == 2 {
			globals.Globals().CaptureFrameFactor = 0.75
		} else if value == 3 {
			globals.Globals().CaptureFrameFactor = 0.5
		}

		this.cameraCaptureFrameFactorWidgetLabel.SetText(fmt.Sprintf("%.2f", globals.Globals().CaptureFrameFactor))
	}

	this.cameraCaptureFPSWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", globals.Globals().CaptureFPS))
	this.cameraCaptureFPSWidget = widget.NewSlider(1, 15)
	this.cameraCaptureFPSWidget.Value = globals.Globals().CaptureFPS
	this.cameraCaptureFPSWidget.OnChangeEnded = func(value float64) {
		globals.Globals().CaptureFPS = value
		this.cameraCaptureFPSWidgetLabel.SetText(fmt.Sprintf("%.0f", globals.Globals().CaptureFPS))
	}

	// ----------------- PaperContour -----------------

	this.paperContourFactorWidgetLabel = widget.NewLabel(fmt.Sprintf("%.2f", globals.Globals().PaperFindContourFactor))
	this.paperContourFactorWidget = widget.NewSlider(1, 10)
	this.paperContourFactorWidget.Value = 1 / globals.Globals().PaperFindContourFactor
	this.paperContourFactorWidget.OnChangeEnded = func(value float64) {
		globals.Globals().PaperFindContourFactor = 1 / value
		this.paperContourFactorWidgetLabel.SetText(fmt.Sprintf("%.2f", globals.Globals().PaperFindContourFactor))
	}

	// ----------------- Circle Detection -----------------

	this.thresholdHoughCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", globals.Globals().ThresholdHoughCircles))
	this.meanFindCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", globals.Globals().MeanFindCircles))
	this.dpHoughCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0f", globals.Globals().DpHoughCircles))
	this.gaussianBlurFindCirclesWidgetLabel = widget.NewLabel(fmt.Sprintf("%.2fmm", globals.Globals().GaussianBlurFindCircles))
	this.adaptiveThresholdBlockSizeWidgetLabel = widget.NewLabel(fmt.Sprintf("%.2fmm", globals.Globals().AdaptiveThresholdBlockSize))
	this.adaptiveThresholdSubtractMeanWidgetLabel = widget.NewLabel(fmt.Sprintf("%.1f", globals.Globals().AdaptiveThresholdSubtractMean))

	this.thresholdHoughCirclesWidget = widget.NewSlider(0, 255)
	this.thresholdHoughCirclesWidget.Value = globals.Globals().ThresholdHoughCircles
	this.thresholdHoughCirclesWidget.OnChangeEnded = func(value float64) {
		globals.Globals().ThresholdHoughCircles = value
		this.thresholdHoughCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}

	this.meanFindCirclesWidget = widget.NewSlider(1, 254)
	this.meanFindCirclesWidget.Value = globals.Globals().MeanFindCircles
	this.meanFindCirclesWidget.OnChangeEnded = func(value float64) {
		globals.Globals().MeanFindCircles = value
		this.meanFindCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}

	this.dpHoughCirclesWidget = widget.NewSlider(0, 3)
	this.dpHoughCirclesWidget.Value = globals.Globals().DpHoughCircles
	this.dpHoughCirclesWidget.OnChangeEnded = func(value float64) {
		globals.Globals().DpHoughCircles = value
		this.dpHoughCirclesWidgetLabel.SetText(fmt.Sprintf("%.0f", value))
	}

	this.gaussianBlurFindCirclesWidget = widget.NewSlider(1, 50)
	this.gaussianBlurFindCirclesWidget.Value = float64(globals.Globals().GaussianBlurFindCircles) * 10
	this.gaussianBlurFindCirclesWidget.OnChangeEnded = func(value float64) {

		globals.Globals().GaussianBlurFindCircles = (value / 10)
		this.gaussianBlurFindCirclesWidgetLabel.SetText(fmt.Sprintf("%.2fmm", globals.Globals().GaussianBlurFindCircles))

	}

	this.adaptiveThresholdBlockSizeWidget = widget.NewSlider(30, 90)
	this.adaptiveThresholdBlockSizeWidget.Value = float64(globals.Globals().AdaptiveThresholdBlockSize) * 10
	this.adaptiveThresholdBlockSizeWidget.OnChangeEnded = func(value float64) {
		globals.Globals().AdaptiveThresholdBlockSize = value / 10
		this.adaptiveThresholdBlockSizeWidgetLabel.SetText(fmt.Sprintf("%.2fmm", (globals.Globals().AdaptiveThresholdBlockSize)))

	}

	this.adaptiveThresholdSubtractMeanWidget = widget.NewSlider(-10, 10)
	this.adaptiveThresholdSubtractMeanWidget.Value = float64(globals.Globals().AdaptiveThresholdSubtractMean)
	this.adaptiveThresholdSubtractMeanWidget.OnChangeEnded = func(value float64) {
		globals.Globals().AdaptiveThresholdSubtractMean = float32(value)
		this.adaptiveThresholdSubtractMeanWidgetLabel.SetText(fmt.Sprintf("%.1f", value))
	}

	this.checkR = widget.NewCheck("R", this.SetFCChannels)
	this.checkR.Checked = (globals.Globals().FindContourChannelMask & 4) != 0
	this.checkG = widget.NewCheck("G", this.SetFCChannels)
	this.checkG.Checked = (globals.Globals().FindContourChannelMask & 2) != 0
	this.checkB = widget.NewCheck("B", this.SetFCChannels)
	this.checkB.Checked = (globals.Globals().FindContourChannelMask & 1) != 0

	this.paperFindContourNoiseBlurSizeWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0fpx", float64(globals.Globals().PaperFindContourNoiseBlurSize)))
	this.paperFindContourNoiseBlurSizeWidget = widget.NewSlider(3, 79)
	this.paperFindContourNoiseBlurSizeWidget.Value = float64(globals.Globals().PaperFindContourNoiseBlurSize)
	this.paperFindContourNoiseBlurSizeWidget.OnChangeEnded = func(value float64) {
		globals.Globals().PaperFindContourNoiseBlurSize = int(value)
		this.paperFindContourNoiseBlurSizeWidgetLabel.SetText(fmt.Sprintf("%.0fpx", float64(globals.Globals().PaperFindContourNoiseBlurSize)))

	}

	this.erodeDillateSizeWidgetLabel = widget.NewLabel(fmt.Sprintf("%.0fpx", float64(globals.Globals().ErodeDillateSize)))
	this.erodeDillateSizeWidget = widget.NewSlider(3, 79)
	this.erodeDillateSizeWidget.Value = float64(globals.Globals().ErodeDillateSize)
	this.erodeDillateSizeWidget.OnChangeEnded = func(value float64) {
		globals.Globals().ErodeDillateSize = int(value)
		this.erodeDillateSizeWidgetLabel.SetText(fmt.Sprintf("%.0fpx", float64(globals.Globals().ErodeDillateSize)))

	}

	//globals.Globals().ShowImage = 502

	outpuList := []struct {
		Name  string
		Index int
	}{
		{"Standard", 0},
		{"Channel Image", 501},
		{"Eroded Last Channel", 502},
		{"Merged", 503},
		{"TesseractROI", 601},
		{"ROIs", 701},
	}
	this.outputWidget = widget.NewSelect([]string{}, func(value string) {
		// fmt.Println("cameraSelectWidget",value)
		for i := 0; i < len(outpuList); i++ {
			if value == outpuList[i].Name {
				globals.Globals().ShowImage = outpuList[i].Index
			}
		}
	})
	for i := 0; i < len(outpuList); i++ {
		this.outputWidget.Options = append(this.outputWidget.Options, outpuList[i].Name)
		if globals.Globals().ShowImage == outpuList[i].Index {
			this.outputWidget.SetSelected(this.outputWidget.Options[i])
		}
	}

	container := container.New(layout.NewVBoxLayout(),

		widget.NewAccordion(
			&widget.AccordionItem{
				Title: "Kamera",
				Detail: container.New(
					layout.NewGridLayout(1),
					widget.NewLabel("Camera"),
					this.cameraSelectWidget,

					widget.NewLabel("Frame Factor"),
					container.NewBorder(nil, nil, nil, this.cameraCaptureFrameFactorWidgetLabel, this.cameraCaptureFrameFactorWidget),

					widget.NewLabel("Output"),
					this.outputWidget,

					/*
						widget.NewLabel("FPS"),
						container.NewBorder( nil, nil, nil, this.cameraCaptureFPSWidgetLabel,this.cameraCaptureFPSWidget ),
					*/
				),
			},
		),

		widget.NewAccordion(
			&widget.AccordionItem{
				Title: "Paper",
				Detail: container.New(
					layout.NewGridLayout(1),
					widget.NewLabel("Contour Factor"),
					container.NewBorder(nil, nil, nil, this.paperContourFactorWidgetLabel, this.paperContourFactorWidget),

					widget.NewLabel("Use Channel"),
					container.NewBorder(nil, nil, nil, widget.NewLabel("RGB"), container.New(layout.NewGridLayout(3), this.checkR, this.checkG, this.checkB)),

					widget.NewLabel("Noise Blur"),
					container.NewBorder(nil, nil, nil, this.paperFindContourNoiseBlurSizeWidgetLabel, this.paperFindContourNoiseBlurSizeWidget),

					widget.NewLabel("ErodeDillate"),
					container.NewBorder(nil, nil, nil, this.erodeDillateSizeWidgetLabel, this.erodeDillateSizeWidget),
				),
			},
		),

		widget.NewAccordion(
			&widget.AccordionItem{
				Title: "Kreisdetetion",
				Detail: container.New(
					layout.NewGridLayout(1),
					widget.NewLabel("Mean Find Circles"),
					container.NewBorder(nil, nil, nil, this.meanFindCirclesWidgetLabel, this.meanFindCirclesWidget),

					widget.NewLabel("Hough Circles Threshold"),
					container.NewBorder(nil, nil, nil, this.thresholdHoughCirclesWidgetLabel, this.thresholdHoughCirclesWidget),

					widget.NewLabel("Inverse ratio of the accumulator"),
					container.NewBorder(nil, nil, nil, this.dpHoughCirclesWidgetLabel, this.dpHoughCirclesWidget),

					widget.NewLabel("Blursize"),
					container.NewBorder(nil, nil, nil, this.gaussianBlurFindCirclesWidgetLabel, this.gaussianBlurFindCirclesWidget),

					widget.NewLabel("Adaptive Threshold Block Size"),
					container.NewBorder(nil, nil, nil, this.adaptiveThresholdBlockSizeWidgetLabel, this.adaptiveThresholdBlockSizeWidget),

					widget.NewLabel("Adaptive Threshold Subtract Mean"),
					container.NewBorder(nil, nil, nil, this.adaptiveThresholdSubtractMeanWidgetLabel, this.adaptiveThresholdSubtractMeanWidget),
				),
			},
		),
	)
	return container
}

func (t *SettingsScreenClass) CreateContainer() *fyne.Container {
	container := container.New(
		layout.NewPaddedLayout(),
		t.makeSettingsForm(),
	)
	return container
}

func NewSettingsScreenClass() *SettingsScreenClass {
	o := &SettingsScreenClass{}
	return o
}
