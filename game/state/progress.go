package state

import "time"

type Progress struct {
	t time.Duration
}

func NewProgress() Progress {
	return Progress{
		t: 0 * time.Second,
	}
}

func (p *Progress) Duration() time.Duration {
	return p.t
}

func (p *Progress) Tick() {
	p.t += 1 * time.Second
}
