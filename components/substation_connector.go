package components

type SubstationSelectorParams struct {
	FromSelector SearchablePickerParams
	ToSelector   SearchablePickerParams
	LineMrid     string
	LineName     string
}
