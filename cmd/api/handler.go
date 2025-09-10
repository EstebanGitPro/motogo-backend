package api

import (

	"github.com/EstebanGitPro/motogo-backend/internal/domain/person"
	
)

type handler struct {
	personService person.Service
}

func New(service person.Service) *handler {
	return &handler{
		personService: service,
	}
}

