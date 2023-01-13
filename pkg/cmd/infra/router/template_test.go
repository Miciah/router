package router

import (
	"github.com/google/go-cmp/cmp"
	templateplugin "github.com/openshift/router/pkg/router/template"
	"testing"
)

func TestPercentDecodingForHeaders(t *testing.T) {
	testCases := []struct {
		description   string
		inputValue    string
		expectedValue []templateplugin.HTTPHeader
		expectError   bool
	}{
		{
			description: "should percent decode the header values",
			inputValue:  "Content-Location%3A%252Fmy-first-blog-post%3ASet",
			expectedValue: []templateplugin.HTTPHeader{
				{Name: "Content-Location", Value: "'/my-first-blog-post'", Action: "Set"},
			},
			expectError: false,
		},
		{
			description: "should percent decode the multiple header values",
			inputValue:  "X-Frame-Options%3ADENY%3ASet%2CX-Cache-Info%3Anot%2Bcacheable%253B%2Bmeta%2Bdata%2Btoo%2Blarge%3ASet%2CX-XSS-Protection%3ADelete%2CX-Source%3A%2525%255Bres.hdr%2528X-Value%2529%252Clower%255D%3ASet",
			expectedValue: []templateplugin.HTTPHeader{
				{Name: "X-Frame-Options", Value: "'DENY'", Action: "Set"},
				{Name: "X-Cache-Info", Value: "'not cacheable; meta data too large'", Action: "Set"},
				{Name: "X-XSS-Protection", Action: "Delete"},
				{Name: "X-Source", Value: "'%[res.hdr(X-Value),lower]'", Action: "Set"},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		switch gotValue, err := parseHeadersToBeSetOrDeleted(tc.inputValue); {
		case !cmp.Equal(gotValue, tc.expectedValue):
			t.Errorf(" expected %s, got %s", tc.expectedValue, gotValue)
		case err == nil && tc.expectError:
			t.Errorf("%s: expected error, got nil", tc.description)
		case err != nil && !tc.expectError:
			t.Errorf("%s: expected success, got error: %v", tc.description, err)
		}
	}
}
