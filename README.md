# OpenSSH ssh_config(5) parser

This package implements a simple parser for the `ssh_config(5)` configuration file used by
OpenSSH.

The current aim is **not** to provide a complete parser for `ssh_config(5)`.

## Example usage

You create `SSHConfig` objects by calling `ReadSSHConfig(filename string)` which can be used to
lookup SSH options for a specific name (which can be either a hostname or a *alias*).

    import (
      "github.com/piger/sshconfig"
    )

    sshConfig, err := ReadSSHConfig("/home/myuser/.ssh/config)
    if err != nil {
      log.Fatal(err)
    }

    hostConfig, err := sshConfig.Lookup("db1.example.com")
    if err != nil {
      log.Fatal(err)
    }

    if opt, ok := hostConfig["hostname"]; ok {
        fmt.Printf("ProxyCommand for db1: %s", opt)
    }

## Notes

- Option names are converted to lowercase during parsing of the configuration file
