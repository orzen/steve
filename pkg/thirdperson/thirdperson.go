package thirdperson

import (
	"github.com/orzen/3rdperson/pkg/cfg"
	"github.com/orzen/3rdperson/pkg/observer"
)

func Run(c *cfg.Cfg) {
	var observers []*observer.Observer
	var isRunning bool = true

	// Setup actions and observers
	for _, in := range c.Observers {
		o := observer.New(in)
		observers = append(observers, o)
	}

	// Mainloop
	// - Check observer status
	for isRunning {
		select {
		case ok := <-c.Passed:
			c.cmd
			// reload command cmd

		}
	}
}
