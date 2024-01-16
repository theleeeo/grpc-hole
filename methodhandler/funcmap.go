package methodhandler

import (
	"text/template"

	"github.com/google/uuid"
)

var (
	funcMap = template.FuncMap{
		"uuid": uuidFunc,
	}
)

func uuidFunc() string {
	return uuid.NewString()
}
