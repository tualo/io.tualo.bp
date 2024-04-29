package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/canvas"
	"image"
	"gocv.io/x/gocv"
	"time"
	// assets "io.tualo.bp/assets"
	globals "io.tualo.bp/globals"
	"log"
)

type MainScreenClass struct {
	playState bool
	button *widget.Button
	globals *globals.GlobalValuesClass
	fullNameWidget *widget.Label
	boxLabelWidget *widget.Label
	stackLabelWidget *widget.Label
	ballotLabelWidget *widget.Label

	displayImage *canvas.Image
	settingsScreenClass *SettingsScreenClass
	settingsContainer *fyne.Container
	showImage bool
	onLogout func()

	channel chan gocv.Mat
	boxBarcode chan string
	stackBarcode chan string
	ballotBarcode chan string

	ticker *time.Ticker
	main fyne.CanvasObject

}

func (t *MainScreenClass) SetOnLogout(onLogout func()) {
	t.onLogout = onLogout
}

func (t *MainScreenClass) SetChannel(channel chan gocv.Mat, boxBarcode chan string, stackBarcode chan string, ballotBarcode chan string) {
	t.channel = channel
	t.boxBarcode = boxBarcode
	t.stackBarcode = stackBarcode
	t.ballotBarcode = ballotBarcode


}

func (t *MainScreenClass) RedrawImage() {
	

	for range t.ticker.C {

		if len(t.channel)>0 {
			img,ok1 := <- t.channel
			if ok1 {
				t.displayImage.Image = t.matToImage(img)
				t.displayImage.Refresh()
				img.Close()
			}
		}

		if len(t.boxBarcode)>0 {
			boxBarcode,ok2 := <- t.boxBarcode
			if ok2 {
				t.boxLabelWidget.SetText("Kiste: "+boxBarcode)
			}
		}

		if len(t.stackBarcode)>0 {
			stackBarcode,ok3 := <- t.stackBarcode
			if ok3 {
				t.stackLabelWidget.SetText("Stapel: "+stackBarcode)
			}
		}

		if len(t.ballotBarcode)>0 {
			ballotBarcode,ok4 := <- t.ballotBarcode
			if ok4 {
				t.ballotLabelWidget.SetText("Stimmzettel: "+ballotBarcode)
			}
		}

	}
	log.Println("RedrawImage exited")
}


func (t *MainScreenClass) SetFullName(name string) {
	t.fullNameWidget.SetText(name)
}

func (t *MainScreenClass) GetPlayState() bool {
	return t.playState
}

func (t *MainScreenClass) SetPlayState(state bool) {
	t.playState = state
	if state {
		t.button.SetText("Stop")
		t.ticker = time.NewTicker(1 * time.Millisecond)
		go t.RedrawImage()
	} else {
		t.ticker.Stop()
		t.button.SetText("Start")
	}
}

func (t *MainScreenClass) matToImage(mat gocv.Mat) image.Image {
	img, _ := mat.ToImage()
	return img
}


func (t *MainScreenClass) makeMain() fyne.CanvasObject {

/*
	
	f, oserr := os.Open("asset/Image.png")
	if oserr != nil {
		fmt.Println(oserr)
		//os.Exit(1)
	}else{
		defer f.Close()

		g, decodeerr := png.Decode(f)
		if decodeerr != nil {
			fmt.Println(decodeerr)
			//os.Exit(1)
		}
		t.displayImage = canvas.NewImageFromImage(g)
	}
	*/
	t.displayImage  =  canvas.NewImageFromImage(t.matToImage(gocv.NewMatWithSize(640, 480, gocv.MatTypeCV8UC3)))
	//canvas.NewImageFromResource(assets.Image())	
	
	t.displayImage.FillMode = canvas.ImageFillContain
	if !t.showImage {
		t.displayImage.Hide()
	}
	return container.NewBorder(
		nil, 
		nil, 
		nil, 
		nil,  
		t.displayImage,
	)
}


func (t *MainScreenClass) makeTopBar() fyne.CanvasObject {
	t.fullNameWidget = widget.NewLabel("Fullname")
	t.boxLabelWidget = widget.NewLabel("Kiste: UNBEKANNT")
	t.stackLabelWidget = widget.NewLabel("Stapel: UNBEKANNT")
	t.ballotLabelWidget = widget.NewLabel("Stimmzettel: UNBEKANNT")
	return container.New(
		layout.NewHBoxLayout(), 
		
		t.boxLabelWidget,
		t.stackLabelWidget,
		t.ballotLabelWidget,
		layout.NewSpacer(), 

		t.fullNameWidget,

		widget.NewButtonWithIcon("", theme.SettingsIcon(), func() {
			if t.settingsContainer.Visible() {
				t.settingsContainer.Hide()
			}else{
				t.settingsContainer.Show()
			}
			t.main.Refresh()
			t.displayImage.Refresh()
		 }),

		widget.NewButtonWithIcon("Logout", theme.LogoutIcon(), func() {
			if t.onLogout != nil {
				t.onLogout()
			}
		}),

		/*
		t.makeToolbarTab(),
		*/
		
	)

	
	
}



func (this *MainScreenClass) SetGlobals(globals *globals.GlobalValuesClass) {
	this.globals = globals
}

func (t *MainScreenClass) makeOuterContainer(onStartStopCamera func()) fyne.CanvasObject {
	t.button = widget.NewButton("Start/Stop", onStartStopCamera)
	t.button.SetText("Start")

	t.settingsScreenClass = NewSettingsScreenClass()
	t.settingsScreenClass.SetGlobals( t.globals )
	t.settingsContainer = t.settingsScreenClass.CreateContainer()
	
	t.settingsContainer.Hide()
	t.main = t.makeMain()
	return container.NewBorder(
		t.makeTopBar( ), 
		t.button, 
		nil, 
		t.settingsContainer,  
		t.main,
	)
}




func (t *MainScreenClass) CreateContainer(onStartStopCamera func()) *fyne.Container {
	container := container.New(
		layout.NewPaddedLayout(),
		t.makeOuterContainer(onStartStopCamera),
	)
	return container
}

func NewMainScreenClass() *MainScreenClass {
	o := &MainScreenClass{
		playState: false,
		showImage: true,
		onLogout: nil,

	}
	// o.SetPlayState( false )
	return o
}