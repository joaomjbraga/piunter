package modules

import (
	"github.com/joaomjbraga/piunter/pkg/types"
)

type Module interface {
	ID() string
	Name() string
	Description() string
	IsAvailable() bool
	Analyze(threshold int) (*types.AnalysisResult, error)
	Clean(dryRun bool) (*types.CleaningResult, error)
}

type BaseModule struct {
	id          string
	name        string
	description string
}

func (m *BaseModule) ID() string         { return m.id }
func (m *BaseModule) Name() string       { return m.name }
func (m *BaseModule) Description() string { return m.description }