package tesseract

import (
	"github.com/otiai10/gosseract/v2"
	"tualo.de/deep-test/globals"
)

type Tesseract struct {
	Client *gosseract.Client
}

var static *Tesseract

func Static() *Tesseract {
	if static == nil {
		static = &Tesseract{}
		static.Client = gosseract.NewClient()
		// defer static.Client.Close()

		if globals.Globals().TesseractPrefix != "" {
			static.Client.SetTessdataPrefix(globals.Globals().TesseractPrefix)
		}
		static.Client.SetLanguage("deu")
	}
	return static
}
