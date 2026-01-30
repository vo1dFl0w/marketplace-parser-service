package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/marketplace-parser-service/internal/adapters/parsers"
)

func TestParsers_ParseStringToFloat64(t *testing.T) {
	testCases := []struct {
		name   string
		numStr string
		expErr bool
	}{
		{
			name:   "parse with ','",
			numStr: "250,0",
			expErr: false,
		},
		{
			name:   "parse with '.'",
			numStr: "250.0",
			expErr: false,
		},
		{
			name:   "invalid",
			numStr: "invalid",
			expErr: true,
		},
		{
			name:   "empty",
			numStr: "",
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expErr {
				res, err := parsers.ParseStringToFloat64(tc.numStr)
				assert.NoError(t, err)
				assert.NotNil(t, res)
			} else {
				_, err := parsers.ParseStringToFloat64(tc.numStr)
				assert.Error(t, err)
			}
		})
	}
}

func TestParsers_ParseStringToInteger(t *testing.T) {
	testCases := []struct {
		name   string
		numStr string
		expErr bool
	}{
		{
			name:   "valid",
			numStr: "250",
			expErr: false,
		},
		{
			name:   "invalid",
			numStr: "invalid",
			expErr: true,
		},
		{
			name:   "empty",
			numStr: "",
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expErr {
				res, err := parsers.ParseStringToInteger(tc.numStr)
				assert.NoError(t, err)
				assert.NotNil(t, res)
			} else {
				_, err := parsers.ParseStringToInteger(tc.numStr)
				assert.Error(t, err)
			}
		})
	}
}