package ui

import (
	"image/color"
	"fyne.io/fyne/v2"
	

	svg "io.tualo.bp/svg"
)


// InvertedThemedResource is a resource wrapper that will return a version of the resource with the main color changed
// for use over highlighted elements.
type SuccessThemedResource struct {
	source fyne.Resource
	color color.RGBA
}

// NewInvertedThemedResource creates a resource that adapts to the current theme for use over highlighted elements.
func NewSuccessThemedResource(orig fyne.Resource, col color.RGBA) *SuccessThemedResource {
	res := &SuccessThemedResource{source: orig, color: col}
	return res
}

// Name returns the underlying resource name (used for caching).
func (res *SuccessThemedResource) Name() string {
	return "inverted_" + res.source.Name()
}

// Content returns the underlying content of the resource adapted to the current background color.
func (res *SuccessThemedResource) Content() []byte {
	return svg.Colorize(res.source.Content(), res.color)
}

// Original returns the underlying resource that this inverted themed resource was adapted from
func (res *SuccessThemedResource) Original() fyne.Resource {
	return res.source
}