package validation

import (
	"net/http"
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors myErrors
}

// New make an empty form and return it
func New(data url.Values) *Form {
	return &Form{
		data,
		myErrors(map[string][]string{}),
	}
}

// RequiredField use for saving the required msg for each field
func (f *Form) RequiredField(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field is required for us!")
		}
	}
}

// MinLength for validating the count of specific field
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, "Minimum length is not provided")
		return false
	}
	return true
}

// Valid use for validation
func (f *Form) Valid() bool {

	return len(f.Errors) == 0
}

// Has check the availability of a specific field
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field is required")
		return false
	}

	return true
}
