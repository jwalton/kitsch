package gitutils

import (
	"bufio"
	"io"
	"strings"
)

type gitconfigBranch struct {
	// Branch is the name of the git branch.
	Branch string
	// Remote is the name of the remote.
	Remote string
	// Merge is the upstream branch (e.g. "refs/heads/master")
	Merge string
}

type gitconfig struct {
	// Branches are configured git branches.
	Branches map[string]*gitconfigBranch
}

func (utils *GitUtils) localConfig() (gitconfig, error) {
	configFile, err := utils.fsys.Open(".git/config")
	if err != nil {
		return gitconfig{}, err
	}

	defer configFile.Close()

	return parseGitConfig(configFile)
}

func parseGitConfig(reader io.Reader) (gitconfig, error) {
	config := gitconfig{
		Branches: map[string]*gitconfigBranch{},
	}
	currentBranch := ""

	// Read the file line by line.
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		// If we're reaching a branch, add each line to it.
		if currentBranch != "" {
			if strings.HasPrefix(line, "[") {
				currentBranch = ""
			} else {
				key, value := parseGitConfigLine(line)
				switch key {
				case "remote":
					config.Branches[currentBranch].Remote = value
				case "merge":
					config.Branches[currentBranch].Merge = value
				}
				continue
			}
		}

		if strings.HasPrefix(line, "[") {
			key, value := parseGitConfigHeader(line)
			switch key {
			case "branch":
				// Start reading a branch.
				currentBranch = value
				if currentBranch != "" {
					config.Branches[currentBranch] = &gitconfigBranch{Branch: value}
				}
			}
		} else {
			// Skip this line.
			continue
		}
	}

	err := scanner.Err()
	if err != nil && err != io.EOF {
		return config, err
	}

	return config, nil
}

// parseGitConfigHeader parses a git config header line (e.g. `[branch "master"]` or "[core]").
func parseGitConfigHeader(line string) (key string, value string) {
	if !strings.HasPrefix(line, "[") || !strings.HasSuffix(line, "]") {
		return "", ""
	}
	line = line[1 : len(line)-1]

	spaceIndex := strings.Index(line, " ")
	if spaceIndex == -1 {
		// If there's no value, the whole thing is the key.
		return strings.TrimSpace(line), ""
	}

	key = strings.TrimSpace(line[:spaceIndex])
	value = strings.TrimSpace(line[spaceIndex+1:])

	// Remove "s from the value.
	if value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
	}
	return key, value
}

// parseGitConfigLine parses a key/value pair from a git config line (e.g. `remote = origin`).
func parseGitConfigLine(line string) (key string, value string) {
	equalsIndex := strings.Index(line, "=")
	if equalsIndex == -1 {
		return "", ""
	}

	key = strings.TrimSpace(line[:equalsIndex])
	value = strings.TrimSpace(line[equalsIndex+1:])
	return key, value
}
