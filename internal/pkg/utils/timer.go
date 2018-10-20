/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import (
	"fmt"
	"strings"
	"time"
)

type Timer interface {
	Start()
	Stop()
	Duration() string
	StartedAt() time.Time
	StoppedAt() time.Time
}

type timer struct {
	start time.Time
	end   time.Time
}

func NewTimer() Timer {
	return &timer{}
}

func (t *timer) Start() {
	t.start = time.Now()
}

func (t *timer) Stop() {
	t.end = time.Now()
}

func (t *timer) Duration() string {
	return CalculateDuration(t.start, t.end)
}

func (t *timer) StartedAt() time.Time {
	return t.start
}

func (t *timer) StoppedAt() time.Time {
	return t.end
}

func CalculateDuration(start, end time.Time) string {

	d := end.Sub(start)

	var text string

	h := int(d.Hours())
	m := int(d.Minutes()) - (h * 60)
	s := int(d.Seconds()) - (int(d.Minutes()) * 60)

	if h > 0 {
		text += fmt.Sprintf("%dh ", h)
	}

	if m > 0 {
		text += fmt.Sprintf("%dm ", m)
	}

	if s > 0 {
		text += fmt.Sprintf("%ds ", s)
	} else {

		durStr := d.String()

		var finalDur string
		var stripLen int
		stripLen = 3
		if strings.HasSuffix(durStr, "ms") || strings.HasSuffix(durStr, "ns") {
			stripLen = 2
		}

		idx := len(durStr) - stripLen
		dur := durStr[:idx]
		notation := durStr[idx:]

		dotIdx := strings.Index(dur, ".")
		if dotIdx > 0 {

			befrDec := dur[:dotIdx]
			finalDur = befrDec + notation
		} else {
			finalDur = durStr
		}

		text += finalDur
	}

	return strings.TrimSpace(text)
}
