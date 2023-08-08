package service

import "github.com/jhump/protoreflect/desc"

type Service struct {
	Name string

	methods []*desc.MethodDescriptor

	descriptor *desc.ServiceDescriptor
}

func New(name string, descriptor *desc.ServiceDescriptor) *Service {
	return &Service{
		Name:       name,
		descriptor: descriptor,
		methods:    descriptor.GetMethods(),
	}
}

func (s *Service) GetDescriptorFile() *desc.FileDescriptor {
	return s.descriptor.GetFile()
}

func (s *Service) GetMethod(methodName string) *desc.MethodDescriptor {
	for _, m := range s.methods {
		if m.GetName() == methodName {
			return m
		}
	}

	return nil
}
