package observer

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/orzen/3rdperson/pkg/reaction"
	"github.com/rs/zerolog/log"
)

type Input struct {
	Path     string
	Events   string
	Reaction *reaction.Input

	passed chan Feedback `yaml:",omitempty"`
	failed chan Feedback `yaml:",omitempty"`
}

type Observer struct {
	Path     string
	Events   []fsnotify.Op
	Reaction *reaction.Reaction

	isWatching bool
	watcher    *fsnotify.Watcher
	passed     chan Feedback
	failed     chan Feedback
}

type Feedback struct {
	Observer *Observer
	Error    error
}

func StringToEvents(str string) []fsnotify.Op {
	events := []fsnotify.Op{}

	split := strings.Split(str, " ")

	for _, n := range split {
		switch {
		case n == "any":
			return []fsnotify.Op{fsnotify.Create,
				fsnotify.Remove, fsnotify.Write,
				fsnotify.Rename, fsnotify.Chmod}
		case n == "create":
			events = append(events, fsnotify.Create)
		case n == "remove":
			events = append(events, fsnotify.Remove)
		case n == "write":
			events = append(events, fsnotify.Write)
		case n == "rename":
			events = append(events, fsnotify.Rename)
		case n == "chmod":
			events = append(events, fsnotify.Chmod)
		default:
			log.Fatal().Msgf("invalid event: %s", n)
		}
	}

	return events
}

func New(in Input) *Observer {
	return &Observer{
		Path:     in.Path,
		Events:   StringToEvents(in.Events),
		Reaction: reaction.New(*in.Reaction),

		passed: in.passed,
		failed: in.failed,
	}
}

func (a *Observer) Close() {
	a.isWatching = false

	if a.watcher != nil {
		a.watcher.Close()
	}
}

func (a *Observer) triggerReaction() {
	// This should _not_ run in a goroutine since it could cause conflicts
	// if the action is not thread-safe
	if err := a.Reaction.Run(); err != nil {
		log.Error().Msgf("action '%s': %v", a.Path, err)
		a.failed <- Feedback{
			Observer: a,
			Error:    err,
		}

	} else {
		a.passed <- Feedback{
			Observer: a,
		}
	}
}

func (a *Observer) matchEvents(event fsnotify.Event) {
	// Trigger a reaction if the event is matching
	for _, n := range a.Events {
		if event.Op&n == n {
			a.triggerReaction()
		}
	}
}

func (a *Observer) Watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal().Msgf("new watcher '%s': %v", err)
	}

	a.watcher = watcher
	a.isWatching = true

	go func() {
		for a.isWatching {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Debug().Msgf("event: %v", event)

				a.matchEvents(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Error().Err(err).Msgf("watcher '%s'", a.Path)
			}
		}
	}()

	err = watcher.Add(a.Path)
	if err != nil {
		log.Fatal().Msgf("watcher add '%s': %v", a.Path, err)
	}
}
