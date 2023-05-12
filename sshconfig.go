// Copyright 2018 Daniel Kertesz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sshconfig

import (
	"bufio"
	"errors"
	"os"
	"os/user"
	"regexp"
	"strings"
)

// BUG(piger): This parser does not currently supports all the TOKENS used in ssh_config(5).

var (
	patternSeparator = regexp.MustCompile("[ ,]")
	optSeparator     = regexp.MustCompile("(?: *= *| +)")
	hostMatch        = regexp.MustCompile("(?i)^host +(.*)")
)

// SSHOptions is a map containing the SSH configuration for a single hostname
type SSHOptions map[string]string

// SSHConfig represent the entire ssh_config file
type SSHConfig struct {
	blocks []*configBlock
}

// addBlock add a SSH configuration "block" to the list of blocks; a configuration block
// is a series of SSH options for one or more patterns.
func (s *SSHConfig) addBlock(block *configBlock) {
	s.blocks = append(s.blocks, block)
}

// Lookup returns the SSH configuration for a given name; some of the TOKENS used by
// ssh_config(5) will be expanded (notably %h, %p and %r).
func (s *SSHConfig) Lookup(name string) (SSHOptions, error) {
	result := make(SSHOptions)

	for _, block := range s.blocks {
		if match, err := matchPatterns(name, block.Patterns); err == nil && match {
			for key, value := range block.Config {
				if _, ok := result[key]; !ok {
					result[key] = value
				}
			}
		} else if err != nil {
			return result, err
		}
	}

	// and now we should expand some options like ProxyCommand replacing format strings
	// like %h and %p with hostname and port, etc.
	specialOptions := []string{"proxycommand", "controlpath"}
	hostname := name
	if h, ok := result["hostname"]; ok {
		hostname = h
	}

	port := "22"
	if p, ok := result["port"]; ok {
		port = p
	}

	localUser, err := user.Current()
	if err != nil {
		return result, err
	}

	remoteUser := localUser.Username
	if u, ok := result["user"]; ok {
		remoteUser = u
	}

	for _, optionName := range specialOptions {
		if value, ok := result[optionName]; ok {
			// Replace some TOKENS

			// %h
			value = strings.Replace(value, "%h", hostname, -1)

			// %p
			value = strings.Replace(value, "%p", port, -1)

			// %r (remote user)
			value = strings.Replace(value, "%r", remoteUser, -1)

			// %u (local user)
			value = strings.Replace(value, "%u", localUser.Username, -1)

			result[optionName] = value
		}
	}

	return result, nil
}

// ReadSSHConfig creates new SSHConfig objects by reading a ssh_config file.
func ReadSSHConfig(filename string) (*SSHConfig, error) {
	sshConfig := SSHConfig{}

	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	var patterns string
	config := make(SSHOptions)

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		matches := hostMatch.FindStringSubmatch(line)
		if len(matches) == 2 {
			// is a new host pattern
			if patterns != "" {
				// a new configBlock can be pushed
				block := newConfigBlock(patterns, config)
				sshConfig.addBlock(block)
			}
			config = make(SSHOptions)
			patterns = matches[1]
		} else {
			if patterns == "" {
				return nil, errors.New("expected to be in an Host block")
			}

			values := optSeparator.Split(line, 2)
			key := strings.ToLower(values[0])
			value := strings.Trim(values[1], " ")
			config[key] = value
		}
	}

	block := newConfigBlock(patterns, config)
	sshConfig.addBlock(block)

	return &sshConfig, nil
}

// configBlock represents the SSH configuration for one or more patterns
type configBlock struct {
	Patterns []string
	Config   SSHOptions
}

func newConfigBlock(patterns string, config SSHOptions) *configBlock {
	var hostPatterns []string
	for _, pattern := range patternSeparator.Split(patterns, -1) {
		pattern = strings.Trim(pattern, " ")
		if len(pattern) > 0 {
			hostPatterns = append(hostPatterns, pattern)
		}
	}
	block := configBlock{hostPatterns, config}
	return &block
}
