package layer

import "tualo.de/deep-test/structs"

type Layer struct {
	drawLayers    map[string]*structs.LayerStruct
	channel       chan structs.LayerStruct
	removeChannel chan structs.LayerRemoveStruct
	enabled       bool
}

var static *Layer

func Static() *Layer {
	if static == nil {
		static = NewLayer()
	}
	return static
}

func NewLayer() *Layer {
	layer := &Layer{
		drawLayers:    make(map[string]*structs.LayerStruct),
		channel:       make(chan structs.LayerStruct, 10),
		removeChannel: make(chan structs.LayerRemoveStruct, 10),
	}
	layer.enabled = true
	go layer.retreive()
	return layer
}
