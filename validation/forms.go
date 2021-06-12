package validation

import (
	"net/http"
	"net/url"
)

type Form struct {
	url.Values
	Errors myErrors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		myErrors(map[string][]string{}),
	}
}

func (f Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}

	return true
}
