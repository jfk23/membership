package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Error errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors{},
	}
}

func (f *Form) Valid() bool {
	return len(f.Error) == 0
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		data := strings.TrimSpace(f.Get(field))
		if data == "" {
			f.Error.Add(field, "Please enter required field!!")
		} 
	}

}

func (f *Form) MinLength(field string, length int) {
	x := f.Get(field)
	if len(x) < length {
		f.Error.Add(field, fmt.Sprintf("Please enter characters more than %d", length))
	}
}

func (f *Form) EmailValidate(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Error.Add(field, "Please enter vaild email address")
	}
}

// Form.has is for checking required value
func (f *Form) Has(field string) bool {
	if strings.TrimSpace(f.Get(field)) != "" {
		
		return true
	} else {
	f.Error.Add(field, "Please enter required field")
	return false}
	// return (f.Get(field) != "")
}