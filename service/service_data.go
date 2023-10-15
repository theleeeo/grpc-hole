package service

// ServiceData is struct uset do save and load services
type serviceData struct {
	// the fully qualified name of the service
	Name string
	// the the file containing the service
	File string
	// the files containing the dependencies of the service
	DependentFiles []string
}
