package grab

import (
	api "io.tualo.bp/api"
	structs "io.tualo.bp/structs"
)


func (this *GrabcameraClass) sendImageItem( boxbarcode string, stackbarcode string, barcode string, id int, marks string,image string ) bool{
	usePiped := true

	if usePiped {
		item:=structs.SendImageQueueItem{
			BoxBarcode: boxbarcode,
			StackBarcode: stackbarcode,
			Barcode: barcode,
			Id: id,
			Marks: marks,
			Image: image,
		}
		this.sendImageQueue <- item
		return true
	}else{
	
		res,err := api.SendReading(
			boxbarcode,
			stackbarcode,
			barcode,
			id,
			marks,
			image,
		)
		
		if err != nil {
			/*
			log.Println("SendReading ERROR",err)
			this.currentState = this.setState("sendError",this.currentState)
			this.setHistoryItem(this.lastBarcode,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
			*/
			return false
		}else{
			if res.Success {
				return true
			}else{
				return false
			}
				/*
				this.doFindCircles = false
				this.currentState = this.setState("sendDone",this.currentState)
				this.setHistoryItem(this.lastBarcode,this.strCurrentBoxBarcode,this.strCurrentStackBarcode,this.currentState)
			}else{
				log.Println("SendReading ERROR",res.Msg)
			}
			*/
		}
	}
}