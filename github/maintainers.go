package github

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	log "gopkg.in/src-d/go-log.v1"
)

type User struct {
	Name   string
	Handle string
	Email  string
}

func GetMaintainers(owner, name string) ([]*User, error) {
	urls := []string{
		fmt.Sprintf(`https://raw.githubusercontent.com/%s/%s/master/MAINTAINERS`, owner, name),
		fmt.Sprintf(`https://raw.githubusercontent.com/%s/%s/master/MAINTAINERS.md`, owner, name),
		fmt.Sprintf(`https://raw.githubusercontent.com/%s/%s/master/maintainers.md`, owner, name),
	}

	for i, url := range urls {

		resp, err := http.DefaultClient.Get(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 404 {
			continue
		}

		if i != 0 {
			log.Warningf("maintainers file found in non-standard URL: %s", url)
		}

		defer resp.Body.Close()

		s := bufio.NewScanner(resp.Body)
		var maintainers []*User
		for s.Scan() {
			line := string(s.Bytes())
			if len(strings.TrimSpace(line)) == 0 {
				log.Warningf("empty line in mantainers file: %s", url)
				continue
			}
			u, err := ParseMaintainer(line)
			if err != nil {
				return maintainers, err
			}
			maintainers = append(maintainers, u)
		}

		if len(maintainers) == 0 {
			log.Warningf("empty maintainers file: %s", url)
		}
		return maintainers, s.Err()
	}

	return nil, nil
}

var maintainerPattern = regexp.MustCompile(`^\s*\*?\s*(.+?)\s*<(\S+?)>\s*\(@?(\S+?)\)\s*$`)

func ParseMaintainer(s string) (*User, error) {
	matches := maintainerPattern.FindStringSubmatch(s)
	if len(matches) == 0 {
		return nil, fmt.Errorf("invalid maintainer format")
	}

	return &User{
		Name:   matches[1],
		Email:  matches[2],
		Handle: matches[3],
	}, nil
}
