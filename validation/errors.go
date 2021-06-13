package validation

type myErrors map[string][]string

// Add append an error msg to myErrors for a specific field
func (e myErrors) Add(field, msg string) {
	e[field] = append(e[field], msg)
}

// Get get the first error of a specific field
func (e myErrors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
