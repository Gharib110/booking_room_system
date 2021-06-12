package validation

type myErrors map[string][]string

func (e myErrors) Add(field, msg string)  {
	e[field] = append(e[field], msg)
}

func (e myErrors) get(field string) string {
	es := e[field]
	if len(es) == 0{
		return ""
	}
	return es[0]
}
