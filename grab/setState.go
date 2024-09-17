
package grab

import (
	"log"
	structs "io.tualo.bp/structs"
)

func (this *GrabcameraClass) setState(name string,oldState structs.ImageProcessorState) structs.ImageProcessorState{
	state := structs.ImageProcessorState{}
	if false {
		log.Println("setState",name,oldState.Name);
	}

	if oldState.Name == "sendDone" && name != "ballotPaperCode"  && name != "noBarcodeFound" && name != "escaped" {
		return oldState
	} 

	state.Name = name

	if len(this.currentStateChannel)==cap(this.currentStateChannel) {
		txt,_:=<-this.currentStateChannel
		if false {
			log.Println("setState",txt)
		}
	}
	this.currentStateChannel <- name



	if (name == "default") {
		state.Red = 0
		state.Green = 0
		state.Blue = 0
	}
	if (name == "findPaperContour") {
		state.Red = 110
		state.Green = 110
		state.Blue = 110
	}

	if (name == "findPaperContourFailed") {
		state.Red = 155
		state.Green = 110
		state.Blue = 110
	}

	if (name == "detectedPaper") {
		state.Red = 200
		state.Green = 110
		state.Blue = 110
	}

	if (name == "findBarcodes") {
		state.Red = 110
		state.Green = 200
		state.Blue = 110
	}

	if (name == "noBarcodeFound") {
		state.Red = 255
		state.Green = 0
		state.Blue = 0
	}

	if (name == "findBoxBarcodes") {
		state.Red = 110
		state.Green = 110
		state.Blue = 200
	}

	if (name == "findStackBarcodes") {
		state.Red = 200
		state.Green = 200
		state.Blue = 110
	}

	if (name == "ballotPaperCode") {
		state.Red = 0
		state.Green = 0
		state.Blue = 255
	}

	if (name == "ballotPaperDetected") {
		state.Red = 255
		state.Green = 255
		state.Blue = 255
	}

	if (name == "ballotPaperNotDetected") {
		state.Red = 100
		state.Green = 255
		state.Blue = 100
	}

	if (name == "ballotPaperMarksAnalysed") {
		state.Red = 255
		state.Green = 100
		state.Blue = 255
	}

	if (name == "doFindCirclesDone") {
		state.Red = 155
		state.Green = 50
		state.Blue = 155
	}

	if (name == "sendError") {
		state.Red = 255
		state.Green = 120
		state.Blue = 120
	}

	if (name == "sendDone") {
		state.Red = 100
		state.Green = 255
		state.Blue = 100
	}

	if (name == "isCorrect") {
		state.Red = 10
		state.Green = 155
		state.Blue = 10
	}


	return state
}
