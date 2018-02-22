# OpenSSH ssh_config(5) parser

This package implements a simple parser for the `ssh_config(5)` configuration file used by
OpenSSH.

This package at the moment only supports a subset of the syntax used by `ssh_config(5)`; the current
implemented features are:

- consider only the first value found for each option
- expand some tokens (like `%h` and `%p`) for specific options (e.g. `ProxyCommand`)
- `Host` pattern matching and patterns negation are supported

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
