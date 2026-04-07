package layer2

import "tualo.de/deep-test/structs"

type Layer2 struct {
	layers map[string]chan structs.LayerStruct
}

var static *Layer2

func Static() *Layer2 {
	if static == nil {
		static = NewLayer2()
	}
	return static
}

func NewLayer2() *Layer2 {
	layer := &Layer2{
		layers: make(map[string]chan structs.LayerStruct),
	}
	go layer.retreive()
	return layer
}
