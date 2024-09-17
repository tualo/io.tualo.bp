package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/canvas"
	"image"
	"image/color"
	"gocv.io/x/gocv"
	"time"

	// assets "io.tualo.bp/assets"
	globals "io.tualo.bp/globals"
	structs "io.tualo.bp/structs"
	api "io.tualo.bp/api"
	
	"log"
	"bytes"
	"io/ioutil"
	//"os"
	"fmt"

	//"fyne.io/fyne/dialog"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/mp3"
	"github.com/gopxl/beep/speaker"
)

type MainScreenClass struct {
	playState bool
	button *widget.Button
	informButton *widget.Button
	globals *globals.GlobalValuesClass
	fullNameWidget *widget.Label
	boxLabelWidget *widget.Label
	stackLabelWidget *widget.Label
	ballotLabelWidget *widget.Label

	list *widget.List
	historyData []structs.HistoryListItem

	ocrLabelWidget *widget.Label
	stateLabelWidget *widget.Label

	displayImage *canvas.Image
	settingsScreenClass *SettingsScreenClass
	settingsContainer *fyne.Container
	showImage bool
	onLogout func()

	channel chan gocv.Mat
	boxBarcode chan string
	stackBarcode chan string
	ballotBarcode chan string
	escapedImage chan bool


	currentStateChannel chan string
	currentOCRChannel chan string
	listItemChannel chan structs.HistoryListItem
	detectedCodesChannel chan structs.DetectedCodes
	sendImageQueue chan structs.SendImageQueueItem



	ticker *time.Ticker
	main fyne.CanvasObject

	alert1 beep.StreamSeekCloser
	alert2 beep.StreamSeekCloser
	alert3 beep.StreamSeekCloser
}

func (t *MainScreenClass) initializeSounds() {
	

}

func (t *MainScreenClass) PlayAlert1() {

	var format beep.Format;


	
	// f, err := os.Open("assets/sms-alert-1-daniel_simon.mp3")
	/*var file *os.File
	f, err := file.Read(resourceSmsAlert1DanielsimonMp3.StaticContent)
	if err != nil {
		log.Fatal(err)
	}
	*/
	var err error

	buf:=bytes.NewBuffer(resourceSmsAlert1DanielsimonMp3.StaticContent)
	clsr := ioutil.NopCloser(buf)
	t.alert1, format, err = mp3.Decode(clsr)
	if err != nil {
		log.Fatal(err)
	}

	if true {
		log.Println("format",format)
	}


	defer t.alert1.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(t.alert1, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func (t *MainScreenClass) SetOnLogout(onLogout func()) {
	t.onLogout = onLogout
}

func (t *MainScreenClass) SetChannel(
	channel chan gocv.Mat, 
	boxBarcode chan string, 
	stackBarcode chan string, 
	ballotBarcode chan string, 
	escapedImage chan bool, 
	currentStateChannel chan string, 
	currentOCRChannel chan string, 
	listItemChannel chan structs.HistoryListItem,
	detectedCodesChannel chan structs.DetectedCodes,
	sendImageQueue chan structs.SendImageQueueItem,
) {
	t.channel = channel
	t.boxBarcode = boxBarcode
	t.stackBarcode = stackBarcode
	t.ballotBarcode = ballotBarcode
	t.escapedImage = escapedImage
	t.currentStateChannel = currentStateChannel
	t.currentOCRChannel = currentOCRChannel
	t.listItemChannel = listItemChannel
	t.detectedCodesChannel = detectedCodesChannel
	t.sendImageQueue = sendImageQueue



}

func (t *MainScreenClass) SendQueuedItems() {

	for range t.ticker.C {
		if len(t.sendImageQueue)>0 {
			item,ok1 := <- t.sendImageQueue
			if ok1 {
				res,err := api.SendReading(
					item.BoxBarcode,
					item.StackBarcode,
					item.Barcode,
					item.Id,
					item.Marks,
					item.Image,
				)
				
				if err != nil {
					log.Println("SendReading ERROR",err)
				}else{
					log.Println("SendReading OK",res.Success, len(t.sendImageQueue))
				}
			}
		}

	}
}

func (t *MainScreenClass) SendDetectedCodes() {
	

	for range t.ticker.C {

		if len(t.detectedCodesChannel)>0 {
			detect,ok1 := <- t.detectedCodesChannel
			if ok1 {
				res,err := api.SendDetectedCodes(
					detect.BoxBarcode,
					detect.StackBarcode,
					detect.Barcode,
				)
				
				if err != nil {
					log.Println("SendDetectedCodes ERROR",err)
				}else{
					log.Println("SendDetectedCodes OK",res.Success)
				}
			}
		}
		

	}
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

		if len(t.currentStateChannel)>0 {
			stateText,ok5 := <- t.currentStateChannel
			if ok5 {
				t.stateLabelWidget.SetText("Zustand: "+stateText)
			}
		}

		if len(t.currentOCRChannel)>0 {
			ocrText,ok6 := <- t.currentOCRChannel
			if ok6 {
				t.ocrLabelWidget.SetText("OCR: "+ocrText)
			}
		}

		if len(t.listItemChannel)>0 {
			histItem,ok7 := <- t.listItemChannel
			if ok7 {

				if histItem.State=="escaped" {
					t.PlayAlert1()
				}
				found:=false
				for i:=0;i<len(t.historyData);i++ {
					if t.historyData[i].Barcode == histItem.Barcode {
						found=true
						t.historyData[i] = histItem
						break
					}
				}
				if !found {
					if len(t.historyData)>0 {
						if (t.historyData[len(t.historyData)-1].State != "sendDone") {
							t.PlayAlert1()
						}
					}
					
					t.historyData = append(t.historyData, histItem)
				}

				if len(t.historyData)>6 {
					t.historyData = t.historyData[1:]
				}



				t.list.Refresh()
				// t.ocrLabelWidget.SetText("OCR: "+histItem)
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
		go t.SendDetectedCodes()
		go t.SendQueuedItems()
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

	t.ocrLabelWidget = widget.NewLabel("OCR: UNBEKANNT")
	t.stateLabelWidget = widget.NewLabel("Zustand: UNBEKANNT")

	return container.New(
		layout.NewHBoxLayout(), 
		
		t.boxLabelWidget,
		t.stackLabelWidget,
		t.ballotLabelWidget,
		t.ocrLabelWidget,
		t.stateLabelWidget,
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


/*
func (t *MainScreenClass) func make(_ fyne.Window) fyne.CanvasObject {
	data := make([]string, 1000)
	for i := range data {
		data[i] = "Test Item " + strconv.Itoa(i)
	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := container.NewHBox(icon, label)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id == 5 || id == 6 {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id] + "\ntaller")
			} else {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id])
		icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	list.Select(125)
	list.SetItemHeight(5, 50)
	list.SetItemHeight(6, 50)

	return container.NewHSplit(list, container.NewCenter(hbox))
}
	*/

func (t *MainScreenClass) makeLeftContainer()  fyne.CanvasObject {
	t.informButton = widget.NewButton("Inform", func() {
		if t.GetPlayState() {
			t.escapedImage <- true
		}
	})
	t.informButton.Importance = widget.DangerImportance


	t.list = widget.NewList(
		func() int {
			return len(t.historyData)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.DocumentIcon()), 
				widget.NewRichTextFromMarkdown("**PAGINATION**\n\nStackCode"),
			)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			// log.Println("update",id,t.historyData[id].State )
			if t.historyData[id].State == "escaped" {
				item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(theme.NewErrorThemedResource(theme.ErrorIcon()))
			}
			if t.historyData[id].State == "isCorrect" {
				
				item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(NewSuccessThemedResource(theme.DocumentIcon(),color.RGBA{0,55,255,255}))
			}
			if t.historyData[id].State == "sendDone" {
				item.(*fyne.Container).Objects[0].(*widget.Icon).SetResource(NewSuccessThemedResource(theme.ConfirmIcon(),color.RGBA{0,255,0,255}))
			}
			item.(*fyne.Container).Objects[0].(*widget.Icon).Refresh()
			md:=fmt.Sprintf("**%s**\n\n%s",t.historyData[id].Barcode,t.historyData[id].StackBarcode)
			item.(*fyne.Container).Objects[1].(*widget.RichText).ParseMarkdown(md)
			//item.(*fyne.Container).Objects[1].(*widget.Label).SetText(t.historyData[id].Barcode)
			//item.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Label).SetText(t.historyData[id].Barcode)
			// item.(*fyne.Container).Objects[1].(*fyne.Container).Objects[1].(*widget.Label).SetText(t.historyData[id].StackBarcode)
		},
	)

	c:=container.NewBorder(
		nil, //top
		t.informButton, // bottom
		nil, // left
		nil,  // right
		t.list, // center
	)

	return c
}

func (t *MainScreenClass) makeOuterContainer(onStartStopCamera func()) fyne.CanvasObject {
	t.button = widget.NewButton("Start/Stop", onStartStopCamera)
	t.button.SetText("Start")





	t.settingsScreenClass = NewSettingsScreenClass()
	t.settingsScreenClass.SetGlobals( t.globals )
	t.settingsContainer = t.settingsScreenClass.CreateContainer()
	
	t.settingsContainer.Hide()
	t.main = t.makeMain()


	

	
	c:=container.NewBorder(
		t.makeTopBar( ), 
		t.button, 
		t.makeLeftContainer(), 
		t.settingsContainer,  
		t.main,
	)

	

	return c
}


func (t *MainScreenClass) OnTypedKey(k *fyne.KeyEvent,onStartStopCamera func()){
	if k.Name == "Space" {
		if t.GetPlayState() {
			//t.SetPlayState(false);
			onStartStopCamera(	);
		}else{
			//t.SetPlayState(true);
			onStartStopCamera(	);
		}
		log.Println(">>>>")
	}
	if k.Name == "Escape" {
		if t.GetPlayState() {
			t.escapedImage <- true
			log.Println("<<<<<<<<<<<<")
		}
	}

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
	/*
	o.initializeSounds()
	o.PlayAlert1()
	*/

	

	// o.SetPlayState( false )
	return o
}