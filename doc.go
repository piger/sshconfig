// Copyright 2018 Daniel Kertesz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
   Package sshconfig provides a simple parser for the ssh_config(5) configuration file
   used by OpenSSH.

   For more informations about the format of the configuration file see:
   https://man.openbsd.org/ssh_config

   Usage

   To read a ssh_config(5) file:
       sshConfig, err := ReadSSHConfig("/home/myuser/.ssh/config)

   To lookup the SSH configuration for a given name:
       hostConfig, err := sshConfig.Lookup("db1.example.com")
       if err != nil {
         log.Fatal(err)
       }

       if opt, ok := hostConfig["hostname"]; ok {
         fmt.Printf("ProxyCommand for db1: %s", opt)
       }
*/
package sshconfig
