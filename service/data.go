package service

import "time"

const (
	ServiceDataFileName = "data.json"
)

// ServiceData is struct uset do save and load services
type serviceData struct {
	// the fully qualified name of the service
	Name string
	// the file containing the service
	File string
	// the files containing the dependencies of the service
	DependentFiles []string
	// The time the serivce was saved
	SavedAt time.Time
}
