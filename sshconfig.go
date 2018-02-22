package sshconfig

import (
	"bufio"
	"log"
	"os"
	"os/user"
	"regexp"
	"strings"
)

var patternSeparator *regexp.Regexp = regexp.MustCompile("[ ,]")
var optSeparator *regexp.Regexp = regexp.MustCompile("(?: *= *| +)")

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

// Lookup return the SSH configuration for a given name
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
			if strings.Index(value, "%h") >= 0 {
				value = strings.Replace(value, "%h", hostname, -1)
			}

			// %p
			if strings.Index(value, "%p") >= 0 {
				value = strings.Replace(value, "%p", port, -1)
			}

			// %r (remote user)
			if strings.Index(value, "%r") >= 0 {
				value = strings.Replace(value, "%r", remoteUser, -1)
			}

			// %u (local user)
			if strings.Index(value, "%u") >= 0 {
				value = strings.Replace(value, "%u", localUser.Username, -1)
			}

			result[optionName] = value
		}
	}

	return result, nil
}

// ReadSSHConfig creates new SSHConfig objects by reading a ssh_config file.
func ReadSSHConfig(filename string) (*SSHConfig, error) {
	sshConfig := SSHConfig{}

	fh, error := os.Open(filename)
	if error != nil {
		log.Fatal(error)
	}
	defer fh.Close()

	var patterns string
	config := make(SSHOptions)
	reHost, err := regexp.Compile("(?i)^host +(.*)")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " ")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		matches := reHost.FindStringSubmatch(line)
		if len(matches) == 2 {
			// is a new host pattern
			if patterns != "" {
				// a new configBlock can be pushed
				block := newConfigBlock(patterns, config)
				sshConfig.addBlock(block)
			}
			patterns = ""
			config = make(SSHOptions)
			patterns = matches[1]
		} else {
			if patterns == "" {
				log.Fatal("Expected to be in a Host block")
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
