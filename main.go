package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2"
	ui "io.tualo.bp/ui"
	theme "io.tualo.bp/ui/theme"
	config "io.tualo.bp/config"
	grab "io.tualo.bp/grab"
	globals "io.tualo.bp/globals"
	api "io.tualo.bp/api"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"log"
)

var topWindow fyne.Window
var loginContainer *fyne.Container
var cameraContainer *fyne.Container
var mainScreenClass *ui.MainScreenClass
var loginScreenClass *ui.LoginScreenClass
var settingsScreenClass *ui.SettingsScreenClass

var configData *config.ConfigurationClass
var g *globals.GlobalValuesClass

var appID = "io.tualo.bp"

func main() {


	configData = config.NewConfigurationClass()
	configData.SetAppID(appID)
	configData.Load()


	g = globals.NewGlobalValuesClass()
	g.SetDefaults()
	log.Println("globals",g)
	g.ConfigData = configData
	g.Load()


	log.Println("globals",g)
	

	grabber := grab.NewGrabcameraClass()
	grabber.SetGlobalValues( g )

	a := app.NewWithID(appID)
	w := a.NewWindow("tualo - ballot scanner")
	topWindow = w
	fyne.CurrentApp().Settings().SetTheme(theme.DefaultTheme())

	loginScreenClass = ui.NewLoginScreenClass()
	loginScreenClass.SetConfig(configData)
	loginScreenClass.SetOnLogin(func(name string) {
		loginContainer.Hide()
		cameraContainer.Show()
		mainScreenClass.SetFullName(name)
	})
		
	loginContainer = loginScreenClass.CreateContainer()

	mainScreenClass = ui.NewMainScreenClass()

	mainScreenClass.SetGlobals(g)



	startStop:=func() {

		if !mainScreenClass.GetPlayState() {
			conf,err := api.GetConfig()

			if err != nil {
				log.Println("GetConfig ERROR",err)
				return
			}
			log.Println("GetConfig",conf)
			log.Println("GetConfig ROIS",conf[0].Rois)
			
			grabber.SetDocumentConfigurations(conf)
		}
		
		mainScreenClass.SetChannel(grabber.GetChannel())
		grabber.SetRun(!mainScreenClass.GetPlayState())
		mainScreenClass.SetPlayState(!mainScreenClass.GetPlayState())
	}
	
	cameraContainer = mainScreenClass.CreateContainer(startStop)

	mainScreenClass.SetOnLogout(func( ) {
		cameraContainer.Hide()
		loginContainer.Show()
	})

	
	content := container.New(
		layout.NewStackLayout(), 
		loginContainer,
		cameraContainer,
		
	)
	loginContainer.Show()
	cameraContainer.Hide()

	w.SetContent(content)

	w.SetMaster()

	w.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		log.Println(k.Name)
		if (cameraContainer.Visible()){
			mainScreenClass.OnTypedKey(k,startStop)
		}
    })

	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}