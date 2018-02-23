// Copyright 2018 Daniel Kertesz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sshconfig

import (
	"path/filepath"
	"strings"
)

// matchPatterns match a name (usually a hostname) against a list of shell-style patterns.
// Patterns can be negated by using the prefix "!".
func matchPatterns(name string, patterns []string) (bool, error) {
	result := false

	for _, pattern := range patterns {
		negate := false
		if strings.HasPrefix(pattern, "!") {
			negate = true
			pattern = pattern[1:]
		}

		if matched, err := filepath.Match(pattern, name); err == nil && matched {
			if negate {
				return false, nil
			}
			result = true
		} else if err != nil {
			return false, err
		}
	}
	return result, nil
}
