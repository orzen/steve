package cfg

import (
	"github.com/orzen/3rdperson/pkg/observer"
	"github.com/orzen/3rdperson/pkg/utils"
)

type Cfg struct {
	Observers  []observer.Input
	LaunchFunc utils.CmdFunc
	ReloadFunc utils.CmdFunc

	Passed chan observer.Feedback `yaml:",omitempty"`
	Failed chan observer.Feedback `yaml:",omitempty"`
}

func New() *Cfg {
	return &Cfg{
		Observers: []observer.Input{},

		Passed: make(chan observer.Feedback),
		Failed: make(chan observer.Feedback),
	}
}
