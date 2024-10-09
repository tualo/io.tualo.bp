package globals

import (
	config "io.tualo.bp/config"
	structs "io.tualo.bp/structs"
)

type GlobalValuesClass struct {
	ConfigData         *config.ConfigurationClass
	IntCamera          int
	CaptureFrameFactor float64
	CaptureFPS         float64

	PaperFindContourFactor float64

	SumMarksAVG                      float64
	RunVideo                         bool
	LogGrabcamera                    bool
	ShowOutputImage                  bool
	ShowPaperImage                   bool
	ShowCirlceImage                  bool
	ShowDebugList                    bool
	InnerOverdrawDrawCircles         int
	OuterOverdrawDrawCircles         int
	MeanFindCircles                  float64
	DpHoughCircles                   float64
	MinDistHoughCircles              float64
	ThresholdHoughCircles            float64
	AccumulatorThresholdHoughCircles float64
	GaussianBlurFindCircles          float64
	AdaptiveThresholdBlockSize       float64
	AdaptiveThresholdSubtractMean    float32
	TesseractPrefix                  string
	ForcedCameraWidth                int
	ForcedCameraHeight               int
	BarcodeScale                     int
	TesseractScale                   int
	ShowOpenCVWindow                 bool
	DocumentConfigurations           structs.DocumentConfigurations

	ShowImage                     int
	FindContourChannelMask        int
	PaperFindContourNoiseBlurSize int
	ErodeDillateSize              int
}

func NewGlobalValuesClass() *GlobalValuesClass {
	o := &GlobalValuesClass{}
	// o.SetPlayState( false )
	return o
}

func (this *GlobalValuesClass) SetDefaults() {
	this.IntCamera = 1

	this.ShowImage = 0
	this.FindContourChannelMask = 7
	this.PaperFindContourNoiseBlurSize = 15

	this.CaptureFrameFactor = 0.5
	this.CaptureFPS = 3.0

	this.PaperFindContourFactor = 0.2

	this.SumMarksAVG = 0.75
	this.RunVideo = false
	this.LogGrabcamera = false
	this.ShowOutputImage = false
	this.ShowPaperImage = true
	this.ShowCirlceImage = false
	this.ShowDebugList = false
	this.InnerOverdrawDrawCircles = 3
	this.OuterOverdrawDrawCircles = 30
	this.MeanFindCircles = 250
	this.DpHoughCircles = 1
	this.MinDistHoughCircles = 50
	this.ThresholdHoughCircles = 90
	this.AccumulatorThresholdHoughCircles = 10
	this.GaussianBlurFindCircles = 1.0
	this.AdaptiveThresholdBlockSize = 9.0
	this.AdaptiveThresholdSubtractMean = 4.0
	this.TesseractPrefix = ""
	this.ForcedCameraWidth = -1
	this.ForcedCameraHeight = -1
	this.BarcodeScale = 1
	this.TesseractScale = 1
	this.ErodeDillateSize = 23
	this.ShowOpenCVWindow = false
}

func (this *GlobalValuesClass) Load() {
	defaults := GlobalValuesClass{}
	defaults.SetDefaults()
	this.SumMarksAVG = this.ConfigData.GetFloat64("settings", "sumMarksAVG", defaults.SumMarksAVG)
	this.RunVideo = this.ConfigData.GetBool("settings", "runVideo", defaults.RunVideo)
	this.ShowOutputImage = this.ConfigData.GetBool("settings", "showOutputImage", defaults.ShowOutputImage)
	this.ShowPaperImage = this.ConfigData.GetBool("settings", "showPaperImage", defaults.ShowPaperImage)
	this.ShowCirlceImage = this.ConfigData.GetBool("settings", "showCirlceImage", defaults.ShowCirlceImage)
	this.ShowDebugList = this.ConfigData.GetBool("settings", "showDebugList", defaults.ShowDebugList)
	this.InnerOverdrawDrawCircles = this.ConfigData.GetInt("settings", "innerOverdrawDrawCircles", defaults.InnerOverdrawDrawCircles)
	this.OuterOverdrawDrawCircles = this.ConfigData.GetInt("settings", "outerOverdrawDrawCircles", defaults.OuterOverdrawDrawCircles)
	this.MeanFindCircles = this.ConfigData.GetFloat64("settings", "meanFindCircles", defaults.MeanFindCircles)
	this.DpHoughCircles = this.ConfigData.GetFloat64("settings", "dpHoughCircles", defaults.DpHoughCircles)
	this.MinDistHoughCircles = this.ConfigData.GetFloat64("settings", "minDistHoughCircles", defaults.MinDistHoughCircles)
	this.ThresholdHoughCircles = this.ConfigData.GetFloat64("settings", "thresholdHoughCircles", defaults.ThresholdHoughCircles)
	this.AccumulatorThresholdHoughCircles = this.ConfigData.GetFloat64("settings", "accumulatorThresholdHoughCircles", defaults.AccumulatorThresholdHoughCircles)
	this.GaussianBlurFindCircles = this.ConfigData.GetFloat64("settings", "gaussianBlurFindCircles", defaults.GaussianBlurFindCircles)

	this.AdaptiveThresholdBlockSize = this.ConfigData.GetFloat64("settings", "adaptiveThresholdBlockSize", defaults.AdaptiveThresholdBlockSize)
	this.AdaptiveThresholdSubtractMean = this.ConfigData.GetFloat32("settings", "adaptiveThresholdSubtractMean", defaults.AdaptiveThresholdSubtractMean)
	this.TesseractPrefix = this.ConfigData.Get("settings", "tesseractPrefix")
	this.ForcedCameraWidth = this.ConfigData.GetInt("settings", "forcedCameraWidth", defaults.ForcedCameraWidth)
	this.ForcedCameraHeight = this.ConfigData.GetInt("settings", "forcedCameraHeight", defaults.ForcedCameraHeight)
	this.BarcodeScale = this.ConfigData.GetInt("settings", "barcodeScale", defaults.BarcodeScale)
	this.TesseractScale = this.ConfigData.GetInt("settings", "tesseractScale", defaults.TesseractScale)
	this.ShowOpenCVWindow = this.ConfigData.GetBool("settings", "showOpenCVWindow", defaults.ShowOpenCVWindow)

	this.IntCamera = this.ConfigData.GetInt("camera", "index", defaults.IntCamera)
	this.CaptureFrameFactor = this.ConfigData.GetFloat64("camera", "captureFrameFactor", defaults.CaptureFrameFactor)
	this.CaptureFPS = this.ConfigData.GetFloat64("camera", "captureFPS", defaults.CaptureFPS)

	this.PaperFindContourFactor = this.ConfigData.GetFloat64("paper", "contourFactor", defaults.PaperFindContourFactor)

	this.ShowImage = this.ConfigData.GetInt("settings", "showImage", defaults.ShowImage)
	this.FindContourChannelMask = this.ConfigData.GetInt("paper", "findContourChannelMask", defaults.FindContourChannelMask)
	this.PaperFindContourNoiseBlurSize = this.ConfigData.GetInt("paper", "paperFindContourNoiseBlurSize", defaults.PaperFindContourNoiseBlurSize)
	this.ErodeDillateSize = this.ConfigData.GetInt("paper", "erodeDillateSize", defaults.ErodeDillateSize)

}

func (this *GlobalValuesClass) Save() {
	this.ConfigData.SetFloat64("settings", "sumMarksAVG", this.SumMarksAVG)
	this.ConfigData.SetBool("settings", "runVideo", this.RunVideo)
	this.ConfigData.SetBool("settings", "showOutputImage", this.ShowOutputImage)
	this.ConfigData.SetBool("settings", "showPaperImage", this.ShowPaperImage)
	this.ConfigData.SetBool("settings", "showCirlceImage", this.ShowCirlceImage)
	this.ConfigData.SetBool("settings", "showDebugList", this.ShowDebugList)
	this.ConfigData.SetInt("settings", "innerOverdrawDrawCircles", this.InnerOverdrawDrawCircles)
	this.ConfigData.SetInt("settings", "outerOverdrawDrawCircles", this.OuterOverdrawDrawCircles)
	this.ConfigData.SetFloat64("settings", "meanFindCircles", this.MeanFindCircles)
	this.ConfigData.SetFloat64("settings", "dpHoughCircles", this.DpHoughCircles)
	this.ConfigData.SetFloat64("settings", "minDistHoughCircles", this.MinDistHoughCircles)
	this.ConfigData.SetFloat64("settings", "thresholdHoughCircles", this.ThresholdHoughCircles)
	this.ConfigData.SetFloat64("settings", "accumulatorThresholdHoughCircles", this.AccumulatorThresholdHoughCircles)
	this.ConfigData.SetFloat64("settings", "gaussianBlurFindCircles", this.GaussianBlurFindCircles)
	this.ConfigData.SetFloat64("settings", "adaptiveThresholdBlockSize", this.AdaptiveThresholdBlockSize)
	this.ConfigData.SetFloat32("settings", "adaptiveThresholdSubtractMean", this.AdaptiveThresholdSubtractMean)
	this.ConfigData.Set("settings", "tesseractPrefix", this.TesseractPrefix)
	this.ConfigData.SetInt("settings", "forcedCameraWidth", this.ForcedCameraWidth)
	this.ConfigData.SetInt("settings", "forcedCameraHeight", this.ForcedCameraHeight)
	this.ConfigData.SetInt("settings", "barcodeScale", this.BarcodeScale)
	this.ConfigData.SetInt("settings", "tesseractScale", this.TesseractScale)
	this.ConfigData.SetBool("settings", "showOpenCVWindow", this.ShowOpenCVWindow)

	this.ConfigData.SetInt("camera", "index", this.IntCamera)
	this.ConfigData.SetFloat64("camera", "captureFrameFactor", this.CaptureFrameFactor)
	this.ConfigData.SetFloat64("camera", "captureFPS", this.CaptureFPS)

	this.ConfigData.SetInt("settings", "showImage", this.ShowImage)

	this.ConfigData.SetInt("paper", "findContourChannelMask", this.FindContourChannelMask)
	this.ConfigData.SetInt("paper", "paperFindContourNoiseBlurSize", this.PaperFindContourNoiseBlurSize)
	this.ConfigData.SetFloat64("paper", "contourFactor", this.PaperFindContourFactor)
	this.ConfigData.SetInt("paper", "erodeDillateSize", this.ErodeDillateSize)

	this.ConfigData.Save()
}
