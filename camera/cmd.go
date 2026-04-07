package camera

import (
	"fmt"
	"log"
	"os"

	"gocv.io/x/gocv"
	"tualo.de/deep-test/args"
	"tualo.de/deep-test/background"
	"tualo.de/deep-test/continuous"
	"tualo.de/deep-test/cv"
	frzr "tualo.de/deep-test/freezer"
	"tualo.de/deep-test/layer"
	"tualo.de/deep-test/layer2"
	ppr "tualo.de/deep-test/paper"
	"tualo.de/deep-test/structs"
)

var windows map[string]*gocv.Window
var px = 0

//var stopFreezer chan bool = make(chan bool, 1)

var showImageChannel chan structs.ShowImageStruct = make(chan structs.ShowImageStruct, 20)

var run bool = false

/*var checkPaperRunning = false
 */
var freeze frzr.Freezer = frzr.Freezer{}
var bckgrnd background.Background = background.Background{}
var paper ppr.Paper = ppr.Paper{}

var lyr = layer.Static()
var lyr2 = layer2.Static()

func GetCameraList() []structs.CameraList {
	cameraList := []structs.CameraList{}
	for i := 0; i < 5; i++ {
		webcam, err := gocv.VideoCaptureDeviceWithAPI(i, 0)
		if err != nil {
			return cameraList
		}
		fmt.Println("Cam: ", i, webcam.Get(gocv.VideoCaptureFrameWidth), webcam.Get(gocv.VideoCaptureFrameHeight))
		cameraList = append(cameraList, structs.CameraList{Width: int(webcam.Get(gocv.VideoCaptureFrameWidth)), Height: int(webcam.Get(gocv.VideoCaptureFrameHeight)), Index: i, Title: fmt.Sprintf("Camera %d", i)})
		webcam.Close()
	}
	return cameraList
}

func show(title string, mat cv.Mat) {
	if windows != nil {
		wnd, prs := windows[title]
		if prs {
			if !mat.GetPointer().Empty() {
				wnd.IMShow(*mat.GetPointer())
				// wnd.ResizeWindow(300, 400)
				if wnd.WaitKey(1) == 27 {
					// run = false
					os.Exit(0)
				}
			}
		} else {
			if !mat.GetPointer().Empty() {
				wnd = gocv.NewWindow(title)
				wnd.IMShow(*mat.GetPointer())
				if title != "final" {
					wnd.ResizeWindow(300, 400)
				} else {
					wnd.ResizeWindow(600, 800)
				}
				wnd.MoveWindow(px, 100)
				px += 300
				if wnd.WaitKey(1) == 27 {
					// run = false
					os.Exit(0)
				}
			}
			windows[title] = wnd
		}
	}
}

func SetRun(run_val bool) {
	run = run_val
	if run == true {
		go Run()
	}
}

func Run() {
	run = true
	arguments := args.Parser()
	if *arguments.EnableCVWindow {

		windows = make(map[string]*gocv.Window)
	}

	webcam, err := gocv.OpenVideoCapture(*arguments.Camera)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", *arguments.Camera)
		return
	}
	defer webcam.Close()

	origianlImg := cv.NewMat("origianlImg")
	defer origianlImg.Close()

	freeze.Init()
	// freeze.SetStopChannel(stopFreezer)
	paper.Init()
	bckgrnd.Init()

	// webcam.Set(gocv.VideoCaptureAutoExposure, 53)

	// run = true
	// go ThresholdFreezer()

	// go CheckFreezer()

	// go CheckPaper()
	// go CheckPaperDebug()
	// go DisplayStats()
	// go CheckBackground()

	continuous.Static().ShowImageChannel = showImageChannel
	continuous.Static().DrawLayer = lyr2
	// isRunning := false
	for run {
		for len(showImageChannel) > 0 {
			data := <-showImageChannel
			show(data.Title, data.Mat)
			data.Mat.Close()
		}

		if ok := webcam.Read(origianlImg.GetPointer()); !ok {
			fmt.Printf("Device closed: %v\n", *arguments.Camera)
			return
		}
		if origianlImg.GetPointer().Empty() {
			continue
		}

		img := cv.NewMat("rotate")
		gocv.Rotate(origianlImg.Get(), img.GetPointer(), gocv.Rotate90CounterClockwise)
		// origianlImg.Close()
		/*
			sample := gocv.IMRead("sample.jpg", gocv.IMReadAnyColor)
			img.Close()
			img = cv.NewFromMat("sample", sample)
		*/

		// lyr.Set("main", img, false, "", 1, 1, 1)
		lyr2.Set("main", img, false, "", 1, 1, 1)

		if !continuous.Static().IsRunning() {
			go continuous.Static().Process(img)
		}

		/*
			imgResult, imgCount := lyr.Draw()
			log.Println("imgCount", imgCount)
			if imgCount > 0 {
				log.Println(imgResult.GetPointer().Size())
				show("final", imgResult)
			}
		*/

		imgResult2, imgCount2 := lyr2.Draw(600, 800)
		log.Println("imgCount2", imgCount2)
		if imgCount2 > 0 {

			show("final2", imgResult2)

			imgResult2.Close()

		}

		img.Close()

	}
}
