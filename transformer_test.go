package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name             string
		term             term
		alphavilleSeries alphavilleSeries
	}{
		{"Transform term to Alphaville Series", term{
			CanonicalName: "Africa Series",
			RawID:         "Nstein_GL_AFTM_GL_164835"},
			alphavilleSeries{
				UUID:      "56a141a4-9894-3559-b25b-d0142f8148ff",
				PrefLabel: "Africa Series",
				AlternativeIdentifiers: alternativeIdentifiers{
					TME:   []string{"TnN0ZWluX0dMX0FGVE1fR0xfMTY0ODM1-U2VyaWVz"},
					Uuids: []string{"56a141a4-9894-3559-b25b-d0142f8148ff"},
				},
				Type: "AlphavilleSeries"}},
	}

	for _, test := range tests {
		expectedAlphavilleSeries := transformAlphavilleSeries(test.term, "Series")
		assert.Equal(test.alphavilleSeries, expectedAlphavilleSeries, fmt.Sprintf("%s: Expected Alphaville Series incorrect", test.name))
	}

}
