package api

import (
	"encoding/json"
	"log"
)

func (me *API) SendDetectedCodes(boxbarcode string, stackbarcode string, barcode string) (KandidatenResponse, error) {
	var response KandidatenResponse
	sb, err := me.Get(me.systemURL + "papervoteoptical/" + boxbarcode + "/" + stackbarcode + "/" + barcode)
	if err != nil {
		log.Println("SendDetectedCodes ERROR", err)
	} else {
		json.Unmarshal([]byte(sb), &response)
	}
	return response, err
}
