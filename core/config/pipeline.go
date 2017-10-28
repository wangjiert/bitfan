package config

import (
	"time"

	fqdn "github.com/ShowMax/go-fqdn"
	uuid "github.com/nu7hatch/gouuid"
)

type PipelineState int

type Pipeline struct {
	ID                 int
	Uuid               string
	Name               string
	Description        string
	ConfigLocation     string
	ConfigHostLocation string

	StartedAt time.Time
	StoppedAt time.Time
}

var pipelineIndex int = 0

func NewPipeline(name, description, configLocation string) *Pipeline {
	pipelineIndex++
	uid, _ := uuid.NewV4()

	return &Pipeline{
		ID:                 pipelineIndex,
		Name:               name,
		Uuid:               uid.String(),
		Description:        description,
		ConfigLocation:     configLocation,
		ConfigHostLocation: fqdn.Get(),
	}
}
