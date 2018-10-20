/*
Copyright 2018 The Elasticshift Authors.
*/
package utils

import (
	"testing"
	"time"
)

func TestCalculateDuration(t *testing.T) {

	start := time.Now()

	dur, _ := time.ParseDuration("1s")
	time.Sleep(dur)

	end := time.Now()
	result := CalculateDuration(start, end)
	if result != "1s" {
		t.Fail()
	}
}

func TestTimer(t *testing.T) {

	timer := NewTimer()
	timer.Start()

	dur, _ := time.ParseDuration("1s")
	time.Sleep(dur)

	timer.Stop()

	result := timer.Duration()

	if result != "1s" {
		t.Fail()
	}
}
