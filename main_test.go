package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSequence(t *testing.T) {
	expected := "===.===...===.===.===...=.===.=...=.=.=...=.......===.=.===.=...===.===.===...===.=.=...="

	m, _ := text2morse("MORSE CODE")
	seq := morse2sequence(m)

	var builder strings.Builder
	for _, sig := range seq {
		if sig {
			builder.WriteRune('=')
		} else {
			builder.WriteRune('.')
		}
	}
	actual := builder.String()

	if actual != expected {
		require.Equal(t, actual, expected)
	}
}

func TestEvent(t *testing.T) {

}
