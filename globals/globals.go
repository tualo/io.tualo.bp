package globals

import (
	config "tualo.de/deep-test/config"
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
	ApproxPolyDPFactor               float64
	TesseractPrefix                  string
	ForcedCameraWidth                int
	ForcedCameraHeight               int
	BarcodeScale                     int
	TesseractScale                   int
	ShowOpenCVWindow                 bool
	// DocumentConfigurations           structs.DocumentConfigurations

	ShowImage                     int
	FindContourChannelMask        int
	PaperFindContourNoiseBlurSize int
	ErodeDillateSize              int
}

var gbl = &GlobalValuesClass{}

func Globals() *GlobalValuesClass {
	return gbl
}

func (me *GlobalValuesClass) SetDefaults() {
	me.IntCamera = 1

	me.ShowImage = 0
	me.FindContourChannelMask = 7
	me.PaperFindContourNoiseBlurSize = 15

	me.CaptureFrameFactor = 0.5
	me.CaptureFPS = 3.0

	me.ApproxPolyDPFactor = 0.01

	me.PaperFindContourFactor = 0.2

	me.SumMarksAVG = 0.75
	me.RunVideo = false
	me.LogGrabcamera = false
	me.ShowOutputImage = false
	me.ShowPaperImage = true
	me.ShowCirlceImage = false
	me.ShowDebugList = false
	me.InnerOverdrawDrawCircles = 3
	me.OuterOverdrawDrawCircles = 30
	me.MeanFindCircles = 250
	me.DpHoughCircles = 1
	me.MinDistHoughCircles = 50
	me.ThresholdHoughCircles = 90
	me.AccumulatorThresholdHoughCircles = 10
	me.GaussianBlurFindCircles = 1.0
	me.AdaptiveThresholdBlockSize = 9.0
	me.AdaptiveThresholdSubtractMean = 4.0
	me.TesseractPrefix = ""
	me.ForcedCameraWidth = -1
	me.ForcedCameraHeight = -1
	me.BarcodeScale = 1
	me.TesseractScale = 1
	me.ErodeDillateSize = 23
	me.ShowOpenCVWindow = false
}

func (me *GlobalValuesClass) Load() {
	defaults := GlobalValuesClass{}
	defaults.SetDefaults()
	me.SumMarksAVG = me.ConfigData.GetFloat64("settings", "sumMarksAVG", defaults.SumMarksAVG)
	me.RunVideo = me.ConfigData.GetBool("settings", "runVideo", defaults.RunVideo)
	me.ShowOutputImage = me.ConfigData.GetBool("settings", "showOutputImage", defaults.ShowOutputImage)
	me.ShowPaperImage = me.ConfigData.GetBool("settings", "showPaperImage", defaults.ShowPaperImage)
	me.ShowCirlceImage = me.ConfigData.GetBool("settings", "showCirlceImage", defaults.ShowCirlceImage)
	me.ShowDebugList = me.ConfigData.GetBool("settings", "showDebugList", defaults.ShowDebugList)
	me.InnerOverdrawDrawCircles = me.ConfigData.GetInt("settings", "innerOverdrawDrawCircles", defaults.InnerOverdrawDrawCircles)
	me.OuterOverdrawDrawCircles = me.ConfigData.GetInt("settings", "outerOverdrawDrawCircles", defaults.OuterOverdrawDrawCircles)
	me.MeanFindCircles = me.ConfigData.GetFloat64("settings", "meanFindCircles", defaults.MeanFindCircles)
	me.DpHoughCircles = me.ConfigData.GetFloat64("settings", "dpHoughCircles", defaults.DpHoughCircles)
	me.MinDistHoughCircles = me.ConfigData.GetFloat64("settings", "minDistHoughCircles", defaults.MinDistHoughCircles)
	me.ThresholdHoughCircles = me.ConfigData.GetFloat64("settings", "thresholdHoughCircles", defaults.ThresholdHoughCircles)
	me.AccumulatorThresholdHoughCircles = me.ConfigData.GetFloat64("settings", "accumulatorThresholdHoughCircles", defaults.AccumulatorThresholdHoughCircles)
	me.GaussianBlurFindCircles = me.ConfigData.GetFloat64("settings", "gaussianBlurFindCircles", defaults.GaussianBlurFindCircles)

	me.AdaptiveThresholdBlockSize = me.ConfigData.GetFloat64("settings", "adaptiveThresholdBlockSize", defaults.AdaptiveThresholdBlockSize)
	me.AdaptiveThresholdSubtractMean = me.ConfigData.GetFloat32("settings", "adaptiveThresholdSubtractMean", defaults.AdaptiveThresholdSubtractMean)
	me.TesseractPrefix = me.ConfigData.Get("settings", "tesseractPrefix")
	me.ForcedCameraWidth = me.ConfigData.GetInt("settings", "forcedCameraWidth", defaults.ForcedCameraWidth)
	me.ForcedCameraHeight = me.ConfigData.GetInt("settings", "forcedCameraHeight", defaults.ForcedCameraHeight)
	me.BarcodeScale = me.ConfigData.GetInt("settings", "barcodeScale", defaults.BarcodeScale)
	me.TesseractScale = me.ConfigData.GetInt("settings", "tesseractScale", defaults.TesseractScale)
	me.ShowOpenCVWindow = me.ConfigData.GetBool("settings", "showOpenCVWindow", defaults.ShowOpenCVWindow)

	me.IntCamera = me.ConfigData.GetInt("camera", "index", defaults.IntCamera)
	me.CaptureFrameFactor = me.ConfigData.GetFloat64("camera", "captureFrameFactor", defaults.CaptureFrameFactor)
	me.CaptureFPS = me.ConfigData.GetFloat64("camera", "captureFPS", defaults.CaptureFPS)

	me.PaperFindContourFactor = me.ConfigData.GetFloat64("paper", "contourFactor", defaults.PaperFindContourFactor)

	me.ShowImage = me.ConfigData.GetInt("settings", "showImage", defaults.ShowImage)
	me.FindContourChannelMask = me.ConfigData.GetInt("paper", "findContourChannelMask", defaults.FindContourChannelMask)
	me.PaperFindContourNoiseBlurSize = me.ConfigData.GetInt("paper", "paperFindContourNoiseBlurSize", defaults.PaperFindContourNoiseBlurSize)
	me.ErodeDillateSize = me.ConfigData.GetInt("paper", "erodeDillateSize", defaults.ErodeDillateSize)

}

func (me *GlobalValuesClass) Save() {
	me.ConfigData.SetFloat64("settings", "sumMarksAVG", me.SumMarksAVG)
	me.ConfigData.SetBool("settings", "runVideo", me.RunVideo)
	me.ConfigData.SetBool("settings", "showOutputImage", me.ShowOutputImage)
	me.ConfigData.SetBool("settings", "showPaperImage", me.ShowPaperImage)
	me.ConfigData.SetBool("settings", "showCirlceImage", me.ShowCirlceImage)
	me.ConfigData.SetBool("settings", "showDebugList", me.ShowDebugList)
	me.ConfigData.SetInt("settings", "innerOverdrawDrawCircles", me.InnerOverdrawDrawCircles)
	me.ConfigData.SetInt("settings", "outerOverdrawDrawCircles", me.OuterOverdrawDrawCircles)
	me.ConfigData.SetFloat64("settings", "meanFindCircles", me.MeanFindCircles)
	me.ConfigData.SetFloat64("settings", "dpHoughCircles", me.DpHoughCircles)
	me.ConfigData.SetFloat64("settings", "minDistHoughCircles", me.MinDistHoughCircles)
	me.ConfigData.SetFloat64("settings", "thresholdHoughCircles", me.ThresholdHoughCircles)
	me.ConfigData.SetFloat64("settings", "accumulatorThresholdHoughCircles", me.AccumulatorThresholdHoughCircles)
	me.ConfigData.SetFloat64("settings", "gaussianBlurFindCircles", me.GaussianBlurFindCircles)
	me.ConfigData.SetFloat64("settings", "adaptiveThresholdBlockSize", me.AdaptiveThresholdBlockSize)
	me.ConfigData.SetFloat32("settings", "adaptiveThresholdSubtractMean", me.AdaptiveThresholdSubtractMean)
	me.ConfigData.Set("settings", "tesseractPrefix", me.TesseractPrefix)
	me.ConfigData.SetInt("settings", "forcedCameraWidth", me.ForcedCameraWidth)
	me.ConfigData.SetInt("settings", "forcedCameraHeight", me.ForcedCameraHeight)
	me.ConfigData.SetInt("settings", "barcodeScale", me.BarcodeScale)
	me.ConfigData.SetInt("settings", "tesseractScale", me.TesseractScale)
	me.ConfigData.SetBool("settings", "showOpenCVWindow", me.ShowOpenCVWindow)

	me.ConfigData.SetInt("camera", "index", me.IntCamera)
	me.ConfigData.SetFloat64("camera", "captureFrameFactor", me.CaptureFrameFactor)
	me.ConfigData.SetFloat64("camera", "captureFPS", me.CaptureFPS)

	me.ConfigData.SetInt("settings", "showImage", me.ShowImage)

	me.ConfigData.SetInt("paper", "findContourChannelMask", me.FindContourChannelMask)
	me.ConfigData.SetInt("paper", "paperFindContourNoiseBlurSize", me.PaperFindContourNoiseBlurSize)
	me.ConfigData.SetFloat64("paper", "contourFactor", me.PaperFindContourFactor)
	me.ConfigData.SetInt("paper", "erodeDillateSize", me.ErodeDillateSize)

	me.ConfigData.Save()
}
