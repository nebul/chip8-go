package main

import (
	"time"
)

type Clock interface {
	Start()
	Stop()
	Tick() <-chan time.Time
}

type FixedClock struct {
	ticker *time.Ticker
	done   chan bool
}

func NewFixedClock() *FixedClock {
	return &FixedClock{
		ticker: time.NewTicker(16 * time.Millisecond), // ~60 Hz
		done:   make(chan bool),
	}
}

func (fixedClock *FixedClock) Start() {
	go func() {
		for {
			select {
			case <-fixedClock.done:
				return
			case <-fixedClock.ticker.C:
				// Tick event, to be handled in emulator's main loop
			}
		}
	}()
}

func (fixedClock *FixedClock) Stop() {
	fixedClock.done <- true
	fixedClock.ticker.Stop()
}

func (fixedClock *FixedClock) Tick() <-chan time.Time {
	return fixedClock.ticker.C
}
