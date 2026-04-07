package layer2

import (
	"fmt"
	"time"

	"tualo.de/deep-test/cv"
	"tualo.de/deep-test/structs"
)

func (me *Layer2) Set(title string, mat cv.Mat, donotshow bool, base string, alpha float64, beta float64, gamma float64) {
	// Sicherheitsüberprüfungen
	if me == nil {
		return
	}
	if me.layers == nil {
		me.layers = make(map[string]chan structs.LayerStruct)
	}

	cloned := mat.Clone(fmt.Sprintf("%s,%s", title, "setlayer"))
	channel, exists := me.layers[title]
	if !exists {
		me.layers[title] = make(chan structs.LayerStruct, 3)
		channel = me.layers[title]
	}
	for len(channel) > 0 {
		oldLayer, ok := <-channel
		//log.Println("KICKING>", title, len(channel))
		if ok {
			oldLayer.Mat.Close()

		}

	}
	//log.Println("SET>", title, len(channel))
	channel <- structs.LayerStruct{
		Title:          title,
		DoNotDraw:      donotshow,
		BasedOn:        base,
		Mat:            cloned,
		Weighted_Alpha: alpha,
		Weighted_Beta:  beta,
		Weighted_Gamma: gamma,
		Timestamp:      time.Now(),
	}
	//log.Println("SET>", title, len(channel), "done")
}
