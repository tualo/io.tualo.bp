package layer

import (
	"fmt"
	"time"

	"tualo.de/deep-test/cv"
	"tualo.de/deep-test/structs"
)

func (me *Layer) Set(title string, mat cv.Mat, donotshow bool, base string, alpha float64, beta float64, gamma float64) {

	if false {
		me.channel <- structs.LayerStruct{
			Title:          title,
			DoNotDraw:      donotshow,
			BasedOn:        base,
			Mat:            mat.Clone(fmt.Sprintf("%s,%s", title, "setlayer")),
			Weighted_Alpha: alpha,
			Weighted_Beta:  beta,
			Weighted_Gamma: gamma,
			Timestamp:      time.Now(),
		}
	}
}
