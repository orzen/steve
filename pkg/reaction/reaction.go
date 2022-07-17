package reaction

import (
	"github.com/orzen/3rdperson/pkg/utils"
	"github.com/rs/zerolog/log"
)

type Input struct {
	Command         []string
	Func            utils.CmdFunc `yaml:",omitempty"`
	RecoveryCommand []string
	RecoveryFunc    utils.CmdFunc `yaml:",omitempty"`
}

type Reaction struct {
	Func         utils.CmdFunc
	RecoveryFunc utils.CmdFunc
}

func New(in Input) *Reaction {
	var f utils.CmdFunc
	var rf utils.CmdFunc

	if len(in.Command) == 0 && in.Func == nil {
		log.Fatal().Msg("Either command or func must be set for action")
	}

	if len(in.Command) != 0 && in.Func != nil {
		log.Fatal().Msg("Command and func set for action")
	}

	if len(in.Command) > 0 {
		f = utils.CmdToFunc(in.Command)
	} else {
		f = in.Func
	}

	if len(in.RecoveryCommand) > 0 {
		rf = utils.CmdToFunc(in.RecoveryCommand)
	} else {
		rf = in.RecoveryFunc
	}

	return &Reaction{
		Func:         f,
		RecoveryFunc: rf,
	}
}

func (a *Reaction) Run() error {
	return a.Func()
}

func (a *Reaction) Recover() error {
	return a.RecoveryFunc()
}
