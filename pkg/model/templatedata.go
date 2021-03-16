package model

type TemplateData struct {
	StringMap map[string]string
	IntMap map[string]int
	FloatMap map[string]float32
	DataMap map[string]interface{}
	CSRFToken string
	FlashMsg string
	WarningMsg string
	ErrorMsg string
}