package ui

import (
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"tualo.de/deep-test/api"
	"tualo.de/deep-test/config"
)

type LoginScreenClass struct {
	loggedIn    bool
	strUrl      string
	strLogin    string
	strPassword string

	onLogin func(name string)

	url      *widget.Entry
	login    *widget.Entry
	password *widget.Entry

	pingResponse       api.PingResponse
	kandidatenResponse api.KandidatenResponse
}

func (o *LoginScreenClass) doLogin() {

	o.strUrl = o.url.Text
	o.strLogin = o.login.Text
	o.strPassword = o.password.Text

	// teste ob die url mit einem slash endet, wenn nicht ergänze es
	if !strings.HasSuffix(o.strUrl, "/") {
		o.strUrl += "/"
	}

	api.Static().SetUrl(o.strUrl)
	loginResponse, err := api.Static().Login(o.strLogin, o.strPassword)

	if err != nil {
		fyne.CurrentApp().SendNotification(&fyne.Notification{
			Title:   "Login failed",
			Content: err.Error(),
		})

	} else {
		if loginResponse.Success {

			api.Static().Switch("bwbriefwahl_muenchen")
			api.Static().GetTitleRegionsConfig()
			api.Static().RoiConfig()
			api.Static().TitleRegions()
			api.Static().CandidateBarcodes()
			api.Static().BallotpaperSizes()

			config.Configuration().Set("credentials", "url", o.strUrl)
			config.Configuration().Set("credentials", "login", o.strLogin)
			config.Configuration().Set("credentials", "password", o.strPassword)
			config.Configuration().Save()

			o.pingResponse, _ = api.Static().Ping()
			// o.kandidatenResponse, _ = api.GetKandidaten()
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

func (o *LoginScreenClass) makeLoginFormTab() fyne.CanvasObject {

	o.url = widget.NewEntry()
	o.url.SetPlaceHolder("URL")

	o.url.SetText(config.Configuration().Get("credentials", "url"))

	o.login = widget.NewEntry()
	o.login.SetPlaceHolder("Benutzername")
	o.login.SetText(config.Configuration().Get("credentials", "login"))
	// email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	o.password = widget.NewPasswordEntry()
	o.password.SetPlaceHolder("Password")
	o.password.SetText(config.Configuration().Get("credentials", "password"))

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

func (o *LoginScreenClass) CreateContainer() *fyne.Container {
	label := canvas.NewText("Anmelden", color.White)
	label.TextSize = 20
	label.Alignment = fyne.TextAlignCenter
	label.TextStyle = fyne.TextStyle{Bold: true}

	loginContainer := container.New(
		layout.NewPaddedLayout(),
		container.New(
			layout.NewVBoxLayout(),
			layout.NewSpacer(),
			label,
			o.makeLoginFormTab(),
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
