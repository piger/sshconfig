package sshconfig

import (
	"testing"
)

func matchOption(config SSHOptions, optname string, optvalue string, t *testing.T) {
	if value, ok := config[optname]; !ok {
		t.Fatalf("Option '%s' was not found in %q", optname, config)
	} else if optvalue != value {
		t.Fatalf("Option '%s' has wrong value: expected '%s', got '%s'", optname, optvalue, value)
	}
}

func TestLookup(t *testing.T) {
	sshConfig, err := ReadSSHConfig("testdata/ssh_config.1")
	if err != nil {
		t.Fatalf("Cannot parse testdata/ssh_config.1: %s", err.Error())
	}

	cfg1, err := sshConfig.Lookup("jumphost")
	if err != nil {
		t.Fatalf("Coult not lookup 'jumphost': %s", err.Error())
	}

	matchOption(cfg1, "hostname", "jump.example.com", t)
	matchOption(cfg1, "user", "proxy", t)
	matchOption(cfg1, "identitiesonly", "on", t)
	matchOption(cfg1, "stricthostkeychecking", "yes", t)

	cfg2, err := sshConfig.Lookup("app1.example.com")
	if err != nil {
		t.Fatalf("Coult not lookup 'app1.example.com': %s", err.Error())
	}

	matchOption(cfg2, "user", "root", t)
	matchOption(cfg2, "proxycommand", "ssh -W app1.example.com:22 jumphost", t)
	matchOption(cfg2, "identitiesonly", "on", t)
	matchOption(cfg2, "identityfile", "~/.ssh/foo", t)
	matchOption(cfg2, "stricthostkeychecking", "no", t)
}
