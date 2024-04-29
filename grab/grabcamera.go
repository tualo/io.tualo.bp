package grab

import (
	"fmt"
	//"os"
    "log"
	"time"
	"image"
	// "sort"
	// "image/color"
	"gocv.io/x/gocv"
	globals "io.tualo.bp/globals"
	structs "io.tualo.bp/structs"
)

type GrabcameraClass struct {
	loadMuster bool
	globals *globals.GlobalValuesClass
	documentConfigurations structs.DocumentConfigurations
	pixelScale float64
	pixelScaleY float64
	runVideo bool
	paperChannelImage chan gocv.Mat
	imageChannelPaper chan gocv.Mat
	currentBoxBarcode chan string
	currentStackBarcode chan string
	ballotBarcode chan string

	/*
	tesseractPrefix string
	intCamera int
	logGrabcamera bool
	sumMarksAVG float64
	forcedCameraWidth int
	forcedCameraHeight int
	tesseractScale int
	dpHoughCircles float64
	minDist float64
	thresholdHoughCircles float64
	accumulatorThresholdHoughCircles float64
	circleSize int
	circleMinDistance int
	gaussianBlurFindCircles int
	adaptiveThresholdBlockSize int
	adaptiveThresholdSubtractMean float32
	meanFindCircles float64
	innerOverdrawDrawCircles int
	outerOverdrawDrawCircles int
	*/

	onNewImageReady func(chan gocv.Mat)
}

func (this *GrabcameraClass) SetDocumentConfigurations( conf structs.DocumentConfigurations) {
	this.documentConfigurations = conf
}


func (this *GrabcameraClass) SetGlobalValues(globals *globals.GlobalValuesClass) {
	this.globals = globals
	
	if this.globals.GaussianBlurFindCircles % 2!=1 {
		this.globals.GaussianBlurFindCircles++
	}

	if this.globals.AdaptiveThresholdBlockSize % 2!=1 {
		this.globals.AdaptiveThresholdBlockSize++
	}

	if this.globals.AdaptiveThresholdBlockSize < 3 {
		this.globals.AdaptiveThresholdBlockSize = 3
	}


}

func (this *GrabcameraClass) GetCameraList() []structs.CameraList {
	cameraList := []structs.CameraList{}
	for i := 0; i < 5; i++ {
		webcam, err := gocv.VideoCaptureDeviceWithAPI(i,0)
		if err != nil {
			return cameraList
		}
		fmt.Println("Cam: ",i, webcam.Get(gocv.VideoCaptureFrameWidth), webcam.Get(gocv.VideoCaptureFrameHeight))
		cameraList = append(cameraList, structs.CameraList{int(webcam.Get(gocv.VideoCaptureFrameWidth)), int(webcam.Get(gocv.VideoCaptureFrameHeight)), i, fmt.Sprintf("Camera %d",i)})
		webcam.Close()
	}
	return cameraList
}


func (this *GrabcameraClass) ResizeMat(img gocv.Mat,width int, height int) gocv.Mat {
	resizeMat := gocv.NewMat()

	if !img.Empty() {
		if img.Cols() >= width && img.Rows() >= height {
			if height>0 && width>0 {
				//fmt.Println("ResizeMat",img.Cols(),img.Rows(),width,height)
				gocv.Resize(img, &resizeMat, image.Point{width, height}, 0, 0, gocv.InterpolationArea)
				img.Close()
			}
		}
	}
	if resizeMat.Empty() {
		return img
	}
	return resizeMat
}

func (this *GrabcameraClass) SetRun( val bool) {
	this.runVideo = val
	if val {
		go this.Grabcamera()
		go this.processImage()
	}
}



func (this *GrabcameraClass) Grabcamera( ) {

	muster := gocv.NewMat()
	var webcam  *gocv.VideoCapture;
	var err error;

	defer muster.Close()
	if this.loadMuster {
		muster = gocv.IMRead("/Users/thomashoffmann/Desktop/muster2.jpg", gocv.IMReadColor)
		//log.Println("grabcamera >>>>>>>>>>>>>>>>>>>>",muster.Cols(),muster.Rows())
		//return
	}else{

	webcam, err = gocv.VideoCaptureDeviceWithAPI(this.globals.IntCamera,0)
	

	if this.globals.LogGrabcamera {
		log.Println("grabcamera >>>>>>>>>>>>>>>>>>>>",this.globals.IntCamera)
	}
	if this.globals.ForcedCameraWidth > 0 {
		webcam.Set(gocv.VideoCaptureFrameWidth, float64(this.globals.ForcedCameraWidth))
	}
	if this.globals.ForcedCameraHeight > 0 {
		webcam.Set(gocv.VideoCaptureFrameHeight, float64(this.globals.ForcedCameraHeight))
	}

	webcam.Set(gocv.VideoCaptureFrameWidth, 4608/2)
	webcam.Set(gocv.VideoCaptureFrameHeight, 3456/2)

	if err != nil {
		fmt.Println("Error opening capture device: ", 0)
		return
	}
	defer webcam.Close()
}

	img := gocv.NewMat()
	if this.loadMuster {
		defer img.Close()
	}
	/*
	checkMarkList := []CheckMarkList{}
	lastReturnType := ReturnType{}
	*/
	for this.runVideo {
		start:=time.Now()
		rotated := gocv.NewMat()

		if this.loadMuster {
			img = muster.Clone()
			rotated= img.Clone()
			img.Close()
		}else{
			webcam.Read(&img)
			gocv.Rotate(img, &rotated, gocv.Rotate90Clockwise)
		}


		if this.globals.LogGrabcamera {
			log.Println("grabcamera >>>>>>>>>>>>>>>>>>>>",rotated.Cols(),rotated.Rows(),time.Since(start),this.runVideo)
		}
		//debug( fmt.Sprintf("grab %s %d %d %d",time.Since(start),rotated.Cols(),rotated.Rows() , os.Getpid() ) )

		/*
		// Videooutput
		if len(cameraChannelImage)==cap(cameraChannelImage) {
			mat,_ := <-cameraChannelImage
			mat.Close()
		}
		cameraCloned := rotated.Clone()
		cameraChannelImage <- cameraCloned
		*/

		// Paper
		if len(this.paperChannelImage)==cap(this.paperChannelImage) {
			mat,_ := <-this.paperChannelImage
			mat.Close()
		}
		paperCloned := rotated.Clone()
		this.paperChannelImage <- paperCloned
		rotated.Close()

		// log.Println("grabcamera >>>>>>>>>>>>>>>>>>>>",len(this.paperChannelImage))
		// this.notifyImage(this.paperChannelImage)
		
	}
	if !this.loadMuster {
		webcam.Close()
	}

}
func (this *GrabcameraClass) GetChannel() (chan gocv.Mat, chan string, chan string, chan string) {
	return this.imageChannelPaper, this.currentBoxBarcode, this.currentStackBarcode, this.ballotBarcode
}

func NewGrabcameraClass() *GrabcameraClass {
	o := &GrabcameraClass{
		globals: nil,
		pixelScale: 1.0,
		pixelScaleY: 1.0,
		runVideo: true,
		paperChannelImage: make(chan gocv.Mat,1),
		imageChannelPaper: make(chan gocv.Mat,1),
		currentBoxBarcode: make(chan string, 1),
		currentStackBarcode: make(chan string, 1),
		ballotBarcode: make(chan string, 1),
	

				/*
			intCamera: 0,
			loadMuster: false,
			logGrabcamera: false,
			pixelScale: 1.0,
			pixelScaleY: 1.0,
			runVideo: true,
			sumMarksAVG: 0.75,
			forcedCameraWidth: 0,
			forcedCameraHeight: 0,
			dpHoughCircles: 1.0,
			minDist: 20.0,
			thresholdHoughCircles: 90.0,
			accumulatorThresholdHoughCircles: 10.0,
			circleSize: 9,
			circleMinDistance: 22,
			gaussianBlurFindCircles: 19,
			adaptiveThresholdBlockSize: 9,
			adaptiveThresholdSubtractMean: 4.0,
			meanFindCircles: 250.0,
			innerOverdrawDrawCircles: 5,
			outerOverdrawDrawCircles: 30,
			tesseractScale: 1,
		*/
	}
	// o.SetPlayState( false )
	return o
}