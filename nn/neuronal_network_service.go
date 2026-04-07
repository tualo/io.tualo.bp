package nn

type NeuronalNetworkService struct {
	URL    string
	Width  int
	Height int
}

var nn = &NeuronalNetworkService{
	URL:    "http://localhost:8080/predict",
	Width:  32,
	Height: 32,
}

func Service() *NeuronalNetworkService {
	return nn
}
