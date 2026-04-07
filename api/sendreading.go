package api

import (
	"encoding/json"
	"fmt"
	"log"
)

func (me *API) SendReading(boxbarcode string, stackbarcode string, barcode string, id int, marks string, image string) (KandidatenResponse, error) {
	var response KandidatenResponse
	data := "boxbarcode=" + boxbarcode + "&stackbarcode=" + stackbarcode + "&barcode=" + barcode + "&id=" + fmt.Sprintf("%d", id) + "&marks=" + marks + "&image=" + image
	// log.Println("SendReading",data)
	sb, err := me.Post(me.systemURL+"papervote/opticaldata", data)

	if err != nil {
		log.Println("SendReading ERROR", err)
	} else {
		json.Unmarshal([]byte(sb), &response)
	}

	return response, err
}
