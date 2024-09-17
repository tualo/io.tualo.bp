package widget

type ListItem struct {
	BaseWidget
	Text      string
	Alignment fyne.TextAlign // The alignment of the text
	Wrapping  fyne.TextWrap  // The wrapping of the text
	TextStyle fyne.TextStyle // The style of the label text

	// The truncation mode of the text
	//
	// Since: 2.4
	Truncation fyne.TextTruncation
	// Importance informs how the label should be styled, i.e. warning or disabled
	//
	// Since: 2.4
	Importance Importance

	provider *RichText
	binder   basicBinder
}

