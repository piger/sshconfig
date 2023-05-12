// Copyright 2018 Daniel Kertesz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sshconfig

import (
	"testing"
)

func TestMatchPatterns(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := matchPatterns(tt.name, tt.patterns)
			if err != nil {
				t.Errorf("matchPatterns(%s, %q): gave error: %s", tt.name, tt.patterns, err.Error())
			}

			if actual != tt.expected {
				t.Errorf("matchPatterns(%s, %q): expected %t, gave %t", tt.name, tt.patterns, tt.expected, actual)
			}
		})
	}
}
