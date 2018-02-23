// Copyright 2018 Daniel Kertesz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sshconfig

import (
	"testing"
)

type testSSHOptions struct {
	sshOptions SSHOptions
	t          *testing.T
}

func (s testSSHOptions) matchOption(optname, optvalue string) {
	if value, ok := s.sshOptions[optname]; !ok {
		s.t.Fatalf("Option '%s' was not found in %q", optname, s.sshOptions)
	} else if optvalue != value {
		s.t.Fatalf("Option '%s' has wrong value: expected '%s', got '%s'", optname, optvalue, value)
	}
}

func readHostConfig(name string, sshConfig *SSHConfig, t *testing.T) testSSHOptions {
	cfg, err := sshConfig.Lookup(name)
	if err != nil {
		t.Fatalf("Cannot lookup '%s': %s", name, err.Error())
	}

	return testSSHOptions{cfg, t}
}

// TestLookup tests the Lookup function for different SSH options and hostnames.
func TestLookup(t *testing.T) {
	sshConfig, err := ReadSSHConfig("testdata/ssh_config.1")
	if err != nil {
		t.Fatalf("Cannot parse testdata/ssh_config.1: %s", err.Error())
	}

	cfg1 := readHostConfig("jumphost", sshConfig, t)
	cfg1.matchOption("hostname", "jump.example.com")
	cfg1.matchOption("user", "proxy")
	cfg1.matchOption("identitiesonly", "yes")
	cfg1.matchOption("stricthostkeychecking", "yes")

	cfg2 := readHostConfig("app1.example.com", sshConfig, t)
	cfg2.matchOption("user", "root")
	cfg2.matchOption("proxycommand", "ssh -W app1.example.com:22 jumphost")
	cfg2.matchOption("identitiesonly", "yes")
	cfg2.matchOption("identityfile", "~/.ssh/foo")
	cfg2.matchOption("stricthostkeychecking", "no")

	cfg3 := readHostConfig("server.internet.org", sshConfig, t)
	cfg3.matchOption("identitiesonly", "yes")
	cfg3.matchOption("stricthostkeychecking", "yes")
}
