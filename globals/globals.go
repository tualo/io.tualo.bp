package globals

import (
	config "io.tualo.bp/config"
	structs "io.tualo.bp/structs"
)


type GlobalValuesClass struct {
	ConfigData *config.ConfigurationClass
	IntCamera int
	SumMarksAVG float64
	RunVideo bool
	LogGrabcamera bool
	ShowOutputImage bool
	ShowPaperImage bool
	ShowCirlceImage bool
	ShowDebugList bool
	InnerOverdrawDrawCircles int
	OuterOverdrawDrawCircles int
	MeanFindCircles float64
	DpHoughCircles float64
	MinDistHoughCircles float64
	ThresholdHoughCircles float64
	AccumulatorThresholdHoughCircles float64
	GaussianBlurFindCircles int
	AdaptiveThresholdBlockSize int
	AdaptiveThresholdSubtractMean float32
	TesseractPrefix string
	ForcedCameraWidth int
	ForcedCameraHeight int
	BarcodeScale int
	TesseractScale int
	ShowOpenCVWindow bool
	DocumentConfigurations structs.DocumentConfigurations
}


func NewGlobalValuesClass() *GlobalValuesClass {
	o := &GlobalValuesClass{}
	// o.SetPlayState( false )
	return o
}

/*
func (o *GlobalValuesClass) SetConfig(cnf *config.ConfigurationClass) {
	o.configData = cnf
}

func (this *GlobalValuesClass) SetFloat32(key string,value float32) {
	switch key {
	case "adaptiveThresholdSubtractMean":
		this.adaptiveThresholdSubtractMean = value
	default:
		return
	}
}

func (this *GlobalValuesClass) SetFloat64(key string,value float64) {
	switch key {
	case "sumMarksAVG":
		this.sumMarksAVG = value
	case "meanFindCircles":
		this.meanFindCircles = value
	case "dpHoughCircles":
		this.dpHoughCircles = value
	case "minDistHoughCircles":
		this.minDistHoughCircles = value
	case "thresholdHoughCircles":
		this.thresholdHoughCircles = value
	case "accumulatorThresholdHoughCircles":
		this.accumulatorThresholdHoughCircles = value
	default:
		return
	}
}

func (this *GlobalValuesClass) SetBool(key string,value bool) {
	switch key {
		case "runVideo":
			this.runVideo = value
		case "showOutputImage":
			this.showOutputImage = value
		case "showPaperImage":
			this.showPaperImage = value
		case "showCirlceImage":
			this.showCirlceImage = value
		case "showDebugList":
			this.showDebugList = value
		case "showOpenCVWindow":
			this.showOpenCVWindow = value
		default:
			return
	}
}

func (this *GlobalValuesClass) GetBool(key string) bool{
	switch key {
		case "runVideo":
			return this.runVideo
		case "showOutputImage":
			return this.showOutputImage
		case "showPaperImage":
			return this.showPaperImage
		case "showCirlceImage":
			return this.showCirlceImage
		case "showDebugList":
			return this.showDebugList
		case "showOpenCVWindow":
			return this.showOpenCVWindow
		default:
			return false
	}
}


func (this *GlobalValuesClass) SetInt(key string,value int) {
	switch key {
		case "innerOverdrawDrawCircles":
			this.innerOverdrawDrawCircles = value
		case "outerOverdrawDrawCircles":
			this.outerOverdrawDrawCircles = value
		case "gaussianBlurFindCircles":
			this.gaussianBlurFindCircles = value
		case "adaptiveThresholdBlockSize":
			this.adaptiveThresholdBlockSize = value
		case "forcedCameraWidth":
			this.forcedCameraWidth = value
		case "forcedCameraHeight":
			this.forcedCameraHeight = value
		case "barcodeScale":
			this.barcodeScale = value
		case "tesseractScale":
			this.tesseractScale = value
		case "intCamera":
			this.intCamera = value
		default:
			return

	}
}

func (this *GlobalValuesClass) GetInt(key string) int{
	
	switch key {
	case "innerOverdrawDrawCircles":
		return this.innerOverdrawDrawCircles
	case "outerOverdrawDrawCircles":
		return this.outerOverdrawDrawCircles
	case "gaussianBlurFindCircles":
		return this.gaussianBlurFindCircles
	case "adaptiveThresholdBlockSize":
		return this.adaptiveThresholdBlockSize
	case "forcedCameraWidth":
		return this.forcedCameraWidth
	case "forcedCameraHeight":
		return this.forcedCameraHeight
	case "barcodeScale":
		return this.barcodeScale
	case "tesseractScale":
		return this.tesseractScale
	case "intCamera":
		return this.intCamera
	default:
		return 0
	}
}

func (this *GlobalValuesClass) GetFloat32(key string) float32{
	switch key {
	case "adaptiveThresholdSubtractMean":
		return this.adaptiveThresholdSubtractMean
	default:
		return 0
	}
}

func (this *GlobalValuesClass) GetFloat64(key string) float64{
	switch key {
	case "sumMarksAVG":
		return this.sumMarksAVG
	case "meanFindCircles":
		return this.meanFindCircles
	case "dpHoughCircles":
		return this.dpHoughCircles
	case "minDistHoughCircles":
		return this.minDistHoughCircles
	case "thresholdHoughCircles":
		return this.thresholdHoughCircles
	case "accumulatorThresholdHoughCircles":
		return this.accumulatorThresholdHoughCircles
	default:
		return 0
	}
}
*/

func (this *GlobalValuesClass) SetDefaults(){
	this.IntCamera = 1
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
	this.GaussianBlurFindCircles = 19
	this.AdaptiveThresholdBlockSize = 9
	this.AdaptiveThresholdSubtractMean = 4.0
	this.TesseractPrefix = ""
	this.ForcedCameraWidth = -1
	this.ForcedCameraHeight = -1
	this.BarcodeScale = 1
	this.TesseractScale = 1
	this.ShowOpenCVWindow = false
}

func (this *GlobalValuesClass) Load() {
	defaults := GlobalValuesClass{}
	defaults.SetDefaults()
	this.SumMarksAVG = this.ConfigData.GetFloat64("settings","sumMarksAVG",defaults.SumMarksAVG)
	this.RunVideo = this.ConfigData.GetBool("settings","runVideo",defaults.RunVideo)
	this.ShowOutputImage = this.ConfigData.GetBool("settings","showOutputImage",defaults.ShowOutputImage)
	this.ShowPaperImage = this.ConfigData.GetBool("settings","showPaperImage",defaults.ShowPaperImage)
	this.ShowCirlceImage = this.ConfigData.GetBool("settings","showCirlceImage",defaults.ShowCirlceImage)
	this.ShowDebugList = this.ConfigData.GetBool("settings","showDebugList",defaults.ShowDebugList)
	this.InnerOverdrawDrawCircles = this.ConfigData.GetInt("settings","innerOverdrawDrawCircles",defaults.InnerOverdrawDrawCircles)
	this.OuterOverdrawDrawCircles = this.ConfigData.GetInt("settings","outerOverdrawDrawCircles",defaults.OuterOverdrawDrawCircles)
	this.MeanFindCircles = this.ConfigData.GetFloat64("settings","meanFindCircles",defaults.MeanFindCircles)
	this.DpHoughCircles = this.ConfigData.GetFloat64("settings","dpHoughCircles",	defaults.DpHoughCircles)
	this.MinDistHoughCircles = this.ConfigData.GetFloat64("settings","minDistHoughCircles",defaults.MinDistHoughCircles)
	this.ThresholdHoughCircles = this.ConfigData.GetFloat64("settings","thresholdHoughCircles",defaults.ThresholdHoughCircles)
	this.AccumulatorThresholdHoughCircles = this.ConfigData.GetFloat64("settings","accumulatorThresholdHoughCircles",defaults.AccumulatorThresholdHoughCircles)
	this.GaussianBlurFindCircles = this.ConfigData.GetInt("settings","gaussianBlurFindCircles",defaults.GaussianBlurFindCircles)
	this.AdaptiveThresholdBlockSize = this.ConfigData.GetInt("settings","adaptiveThresholdBlockSize",defaults.AdaptiveThresholdBlockSize)
	this.AdaptiveThresholdSubtractMean = this.ConfigData.GetFloat32("settings","adaptiveThresholdSubtractMean",defaults.AdaptiveThresholdSubtractMean)
	this.TesseractPrefix = this.ConfigData.Get("settings","tesseractPrefix")
	this.ForcedCameraWidth = this.ConfigData.GetInt("settings","forcedCameraWidth",defaults.ForcedCameraWidth)
	this.ForcedCameraHeight = this.ConfigData.GetInt("settings","forcedCameraHeight",defaults.ForcedCameraHeight)
	this.BarcodeScale = this.ConfigData.GetInt("settings","barcodeScale",defaults.BarcodeScale)
	this.TesseractScale = this.ConfigData.GetInt("settings","tesseractScale",defaults.TesseractScale)
	this.ShowOpenCVWindow = this.ConfigData.GetBool("settings","showOpenCVWindow",defaults.ShowOpenCVWindow)
}

func (this *GlobalValuesClass) Save() {
	this.ConfigData.SetFloat64("settings","sumMarksAVG",this.SumMarksAVG)
	this.ConfigData.SetBool("settings","runVideo",this.RunVideo)
	this.ConfigData.SetBool("settings","showOutputImage",this.ShowOutputImage)
	this.ConfigData.SetBool("settings","showPaperImage",this.ShowPaperImage)
	this.ConfigData.SetBool("settings","showCirlceImage",this.ShowCirlceImage)
	this.ConfigData.SetBool("settings","showDebugList",this.ShowDebugList)
	this.ConfigData.SetInt("settings","innerOverdrawDrawCircles",this.InnerOverdrawDrawCircles)
	this.ConfigData.SetInt("settings","outerOverdrawDrawCircles",this.OuterOverdrawDrawCircles)
	this.ConfigData.SetFloat64("settings","meanFindCircles",this.MeanFindCircles)
	this.ConfigData.SetFloat64("settings","dpHoughCircles",this.DpHoughCircles)
	this.ConfigData.SetFloat64("settings","minDistHoughCircles",this.MinDistHoughCircles)
	this.ConfigData.SetFloat64("settings","thresholdHoughCircles",this.ThresholdHoughCircles)
	this.ConfigData.SetFloat64("settings","accumulatorThresholdHoughCircles",this.AccumulatorThresholdHoughCircles)
	this.ConfigData.SetInt("settings","gaussianBlurFindCircles",this.GaussianBlurFindCircles)
	this.ConfigData.SetInt("settings","adaptiveThresholdBlockSize",this.AdaptiveThresholdBlockSize)
	this.ConfigData.SetFloat32("settings","adaptiveThresholdSubtractMean",this.AdaptiveThresholdSubtractMean)
	this.ConfigData.Set("settings","tesseractPrefix",this.TesseractPrefix)
	this.ConfigData.SetInt("settings","forcedCameraWidth",this.ForcedCameraWidth)
	this.ConfigData.SetInt("settings","forcedCameraHeight",this.ForcedCameraHeight)
	this.ConfigData.SetInt("settings","barcodeScale",this.BarcodeScale)
	this.ConfigData.SetInt("settings","tesseractScale",this.TesseractScale)
	this.ConfigData.SetBool("settings","showOpenCVWindow",this.ShowOpenCVWindow)
	this.ConfigData.Save()
}