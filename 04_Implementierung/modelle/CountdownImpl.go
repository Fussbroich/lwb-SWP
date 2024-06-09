package modelle

import "time"

type countd struct {
	angehalten bool
	restzeit   time.Duration
}

func NewCountdown(t time.Duration) *countd {
	return &countd{restzeit: t}
}

func (c *countd) Setze(t time.Duration) {
	c.restzeit = t
}

func (c *countd) GibRestzeit() time.Duration {
	return c.restzeit
}

func (c *countd) ZieheAb(d time.Duration) {
	if !c.angehalten {
		c.restzeit -= d
		if c.restzeit <= 0 {
			c.restzeit = 0
		}
	}
}

func (c *countd) IstAbgelaufen() bool {
	return c.restzeit <= 0
}

func (c *countd) Halt() {
	c.angehalten = true
}

func (c *countd) Weiter() {
	c.angehalten = false
}
