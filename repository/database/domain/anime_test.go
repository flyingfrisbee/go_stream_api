package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEpisodesAsSliceInterface(t *testing.T) {
	a := Anime{
		Episodes: []Episode{
			{Text: "1", Endpoint: "/some/ep-1"},
			{Text: "2", Endpoint: "/some/ep-2"},
			{Text: "3", Endpoint: "/some/ep-3"},
			{Text: "4", Endpoint: "/some/ep-4"},
		},
	}
	results := a.GetEpisodesAsSliceInterface()

	for i := 0; i < len(results); i++ {
		assert.Equal(t, a.Episodes[i], results[i].(Episode))

	}
}
