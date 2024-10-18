package grab

import (
	"image/color"
	"log"

	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) setHistoryItem(barcode string, boxcode string, stackcode string, currentState structs.ImageProcessorState) {
	histItem := structs.HistoryListItem{
		Barcode:      barcode,
		BoxBarcode:   boxcode,
		StackBarcode: stackcode,
		State:        currentState.Name,
		StateColor:   color.RGBA{uint8(currentState.Red), uint8(currentState.Green), uint8(currentState.Blue), 120},
	}

	if true {
		log.Println("listItemChannel", len(this.listItemChannel))
	}

	if len(this.listItemChannel) == cap(this.listItemChannel) {
		oldItem, _ := <-this.listItemChannel
		if false {
			log.Println("setState oldItem", oldItem)
		}
	}
	this.listItemChannel <- histItem
}
