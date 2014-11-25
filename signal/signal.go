package signal

import (
	"github.com/dlutxx/goblin/utils"
)

type Signal struct {
	listeners []func(utils.Dict)
}

// argNames indicate what args will be included
// when the signal is sent, purely for documentational
// purposes.
func New(argNames ...string) *Signal {
	return &Signal{
		listeners: make([]func(utils.Dict), 0),
	}
}

func (s *Signal) Connect(lsn func(utils.Dict)) {
	s.listeners = append(s.listeners, lsn)
}

func (s *Signal) Send(args utils.Dict) {
	for _, fn := range s.listeners {
		fn(args)
	}
}
