package nn

func (me *NeuronalNetworkService) Setup(url string, width int, height int) {
	me.URL = url
	me.Width = width
	me.Height = height
}

func (me *NeuronalNetworkService) GetURL() string {
	return me.URL
}
