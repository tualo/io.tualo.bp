package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

var appID = "io.tualo.bp"
var topWindow fyne.Window
var loginContainer *fyne.Container
var cameraContainer *fyne.Container
var mainScreenClass *MainScreenClass
var loginScreenClass *LoginScreenClass
var settingsScreenClass *SettingsScreenClass

func StartAndRun() {

	a := app.NewWithID(appID)
	w := a.NewWindow("tualo - ballot scanner")
	topWindow = w
	//	fyne.CurrentApp().Settings().SetTheme(theme.DefaultTheme())

	loginScreenClass = NewLoginScreenClass()
	// loginScreenClass.SetConfig(configData)
	loginScreenClass.SetOnLogin(func(name string) {
		loginContainer.Hide()
		cameraContainer.Show()
		mainScreenClass.SetFullName(name)
	})

	loginContainer = loginScreenClass.CreateContainer()

	mainScreenClass = NewMainScreenClass()
	cameraContainer = mainScreenClass.CreateContainer()

	// mainScreenClass.SetGlobals(g)
	mainScreenClass.TopWindow = topWindow

	mainScreenClass.SetOnLogout(func() {
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
		if cameraContainer.Visible() {
			// mainScreenClass.OnTypedKey(k, startStop)
		}
	})

	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()

}
