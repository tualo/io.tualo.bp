package grab

import (
	api "io.tualo.bp/api"
	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) sendImageItem(boxbarcode string, stackbarcode string, barcode string, id int, marks string, image string) bool {
	usePiped := true

	if usePiped {
		item := structs.SendImageQueueItem{
			BoxBarcode:   boxbarcode,
			StackBarcode: stackbarcode,
			Barcode:      barcode,
			Id:           id,
			Marks:        marks,
			Image:        image,
		}
		this.sendImageQueue <- item
		return true
	} else {

		res, err := api.SendReading(
			boxbarcode,
			stackbarcode,
			barcode,
			id,
			marks,
			image,
		)

		if err != nil {
			return false
		} else {
			if res.Success {
				return true
			} else {
				return false
			}
		}
	}
}
