package model

import "github.com/jfk23/gobookings/internal/forms"

type TemplateData struct {
	StringMap map[string]string
	IntMap map[string]int
	FloatMap map[string]float32
	DataMap map[string]interface{}
	CSRFToken string
	FlashMsg string
	WarningMsg string
	ErrorMsg string
	Form *forms.Form
	IsAuthenticated int
}