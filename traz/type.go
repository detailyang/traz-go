package traz

import (
	"time"

	"github.com/detailyang/size"
)

type Size struct {
	size.Unit
}

func (s *Size) UnmarshalText(text []byte) error {
	var err error
	s.Unit, err = size.Parse(string(text))
	return err
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
