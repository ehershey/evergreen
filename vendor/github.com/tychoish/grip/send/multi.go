package send

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tychoish/grip/level"
	"github.com/tychoish/grip/message"
)

type multiSender struct {
	senders []Sender
	*base
}

// NewMultiSender configures a new sender implementation that takes a
// slice of Sender implementations that dispatches all messages to all
// implementations. This constructor forces all member Senders to have
// the same name and Level configuration. Use NewConfiguredMultiSender
// to construct a similar Sender
//
// Use the AddToMulti helper to add additioanl senders to one of these
// multi Sender implementations after construction.
func NewMultiSender(name string, l LevelInfo, senders []Sender) (Sender, error) {
	if !l.Valid() {
		return nil, fmt.Errorf("invalid level specification: %+v", l)
	}

	if len(senders) == 0 {
		return nil, errors.New("must specify at least one sender when creating a multi sender")
	}

	for _, sender := range senders {
		sender.SetName(name)
		_ = sender.SetLevel(l)
	}

	return &multiSender{senders: senders, base: newBase(name)}, nil
}

// NewConfiguredMultiSender returns a multi sender implementation with
// Sender members, but does not force the senders to have conforming
// name or level values. Use NewMultiSender to construct a list of
// senders with consistent names and level configurations.
//
// Use the AddToMulti helper to add additioanl senders to one of these
// multi Sender implementations after construction.
func NewConfiguredMultiSender(senders ...Sender) Sender {
	s := &multiSender{senders: senders, base: newBase("")}
	_ = s.base.SetLevel(LevelInfo{Default: level.Invalid, Threshold: level.Invalid})

	return s
}

// AddToMulti is a helper function that takes two Sender instances,
// the first of which must be a multi sender. If this is true, then
// AddToMulti adds the second Sender to the first Sender's list of
// Senders.
//
// Returns an error if the first instance is not a multi sender.
func AddToMulti(multi Sender, s Sender) error {
	sender, ok := multi.(*multiSender)
	if !ok {
		return fmt.Errorf("%s is not a multi sender", multi.Name())
	}

	return sender.add(s)
}

func (s *multiSender) Close() error {
	errs := []string{}
	for _, sender := range s.senders {
		if err := sender.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (s *multiSender) add(sender Sender) error {
	sender.SetName(s.base.Name())

	// ignore the error here; if the base value on the multiSender
	// is not set, then senders should just have their own level values.
	_ = sender.SetLevel(s.base.Level())

	s.senders = append(s.senders, sender)
	return nil
}

func (s *multiSender) Name() string { return s.base.Name() }
func (s *multiSender) SetName(n string) {
	s.base.SetName(n)

	for _, sender := range s.senders {
		sender.SetName(n)
	}
}

func (s *multiSender) Level() LevelInfo { return s.base.Level() }
func (s *multiSender) SetLevel(l LevelInfo) error {
	if err := s.base.SetLevel(l); err != nil {
		return err
	}

	for _, sender := range s.senders {
		_ = sender.SetLevel(l)
	}

	return nil
}

func (s *multiSender) Send(m message.Composer) {
	// if the base level isn't valid, then we should let each
	// sender decide for itself, rather than short circuiting here
	bl := s.base.Level()
	if bl.Valid() && !bl.ShouldLog(m) {
		return
	}

	for _, sender := range s.senders {
		sender.Send(m)
	}
}
