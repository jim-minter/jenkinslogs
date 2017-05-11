package detectors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

const (
	ansiblePlaytaskPlayStartPrefix           = "PLAY"
	ansiblePlaytaskTaskStartPrefix           = "TASK"
	ansiblePlaytaskRunningHandlerStartPrefix = "RUNNING HANDLER"
	ansiblePlaytaskFatalEndPrefix            = "fatal: "
)

var ansiblePlaytaskStartRx = regexp.MustCompile(`^((?:PLAY)|(?:TASK)|(?:RUNNING HANDLER)) \[([^\]]*)\] \*+$`)

func AnsiblePlaytaskFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			if !strings.HasPrefix(s[c1], ansiblePlaytaskPlayStartPrefix) &&
				!strings.HasPrefix(s[c1], ansiblePlaytaskTaskStartPrefix) &&
				!strings.HasPrefix(s[c1], ansiblePlaytaskRunningHandlerStartPrefix) {
			}
			s1 = ansiblePlaytaskStartRx.FindStringSubmatch(s[c1])
			if s1 != nil {
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			if s[c2] == "" {
				s2 = []string{s[c2], ""}
				if strings.HasPrefix(s[c2-1], ansiblePlaytaskFatalEndPrefix) {
					s2 = []string{s[c2], ansiblePlaytaskFatalEndPrefix}
				}
				c2++
				break
			}
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, fmt.Sprintf("ansible %s: %s", strings.ToLower(s1[1]), s1[2]), s2[1] == ansiblePlaytaskFatalEndPrefix))
		}

		i = c2
	}

	return groups
}
