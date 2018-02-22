package sshconfig

import (
	"testing"
)

var patternTests = []struct {
	name     string
	patterns []string
	expected bool
}{
	{"example.com", []string{"example.com"}, true},
	{"example.com", []string{"*.com"}, true},
	{"app.example.com", []string{"*.example.com"}, true},
	{"app.example.com", []string{"*.example.com", "!app.example.com"}, false},
	{"foo.bar", []string{"*.com"}, false},
}

func TestMatchPatterns(t *testing.T) {
	for _, pt := range patternTests {
		actual, err := matchPatterns(pt.name, pt.patterns)
		if err != nil {
			t.Errorf("matchPatterns(%s, %q): gave error: %s", pt.name, pt.patterns, err.Error())
		}

		if actual != pt.expected {
			t.Errorf("matchPatterns(%s, %q): expected %t, gave %t", pt.name, pt.patterns, pt.expected, actual)
		}
	}
}
