package templateparse

import "text/template"

var (
	funcMap = template.FuncMap{
		"uuid": uuidFunc,
	}
)
