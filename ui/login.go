package ui

import (
	"os"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	api "io.tualo.bp/api"
	config "io.tualo.bp/config"
)

type LoginScreenClass struct {
	loggedIn bool
	strUrl string
	strLogin string
	strPassword string

	onLogin func(name string)
	configData *config.ConfigurationClass

	url *widget.Entry
	login *widget.Entry
	password *widget.Entry


	pingResponse api.PingResponse
	kandidatenResponse api.KandidatenResponse


}

func (o *LoginScreenClass) doLogin() {

	
	
	o.strUrl = o.url.Text
	o.strLogin = o.login.Text
	o.strPassword = o.password.Text
	/*
	if o.strUrl == "" {
		o.strUrl = strSystemUrl
	}
	if o.strLogin == "" {
		o.strLogin = strSystemLogin
	}
	if o.strPassword == "" {
		o.strPassword = o.strSystemPassword
	}
	*/
	loginResponse, err := api.Login(o.strUrl, o.strLogin, o.strPassword)
	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Login failed",
			Content: err.Error(),
		})
		
	} else {
		if loginResponse.Success {
			
			o.configData.Set("credentials","url",o.strUrl)
			o.configData.Set("credentials","login",o.strLogin)
			o.configData.Set("credentials","password",o.strPassword)
			o.configData.Save()

			/*
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Login successful",
				Content: "Welcome " + loginResponse.Fullname,
			})
			*/
			api.SetSystemURL(o.strUrl)

			o.pingResponse, _ = api.Ping()
			

			o.kandidatenResponse, _ = api.GetKandidaten()
			// fmt.Println(o.kandidatenResponse)
			
			if o.onLogin != nil {
				o.onLogin(loginResponse.Fullname)
			}

		} else {
			/*
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Login failed",
				Content: loginResponse.Msg,
			})
			*/
		}
	}
	
}

func (o *LoginScreenClass) makeLoginFormTab( ) fyne.CanvasObject {

	o.url = widget.NewEntry()
	o.url.SetPlaceHolder("URL")
	o.url.SetText(o.configData.Get("credentials","url"))

	o.login = widget.NewEntry()
	o.login.SetPlaceHolder("Benutzername")
	o.login.SetText(o.configData.Get("credentials","login"))
	// email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	o.password = widget.NewPasswordEntry()
	o.password.SetPlaceHolder("Password")
	o.password.SetText(o.configData.Get("credentials","password"))

	form := &widget.Form{
		SubmitText: "Anmelden",
		CancelText: "Abbrechen",
		Items: []*widget.FormItem{
			{Text: "URL", Widget: o.url, HintText: "Bitte gib die vollständige URL ein."},
			{Text: "Benutzername", Widget: o.login, HintText: "Bitte gib deinen Benutzernamen ein."},
			{Text: "Passwort", Widget: o.password, HintText: "Bitte gib dein Passwort ein."},
		},
		OnCancel: func() {
			os.Exit(0)
		},
		OnSubmit: o.doLogin,
	}
	
	return form
}

func (o *LoginScreenClass) SetOnLogin(fn func(name string)) {
	o.onLogin = fn
}

func (o *LoginScreenClass) SetConfig(cnf *config.ConfigurationClass) {
	o.configData = cnf
}

func (o *LoginScreenClass) CreateContainer( ) *fyne.Container {
	label := canvas.NewText("Anmelden", color.White)
	label.TextSize = 20
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle= fyne.TextStyle{Bold: true}

	loginContainer := container.New(
			layout.NewPaddedLayout(),
			container.New(
			layout.NewVBoxLayout(), 
			layout.NewSpacer(),
			label,
			o.makeLoginFormTab( ),
			layout.NewSpacer(),
		),
	)
	return loginContainer
}

func NewLoginScreenClass() *LoginScreenClass {
	o := &LoginScreenClass{}
	// o.SetPlayState( false )
	return o
}