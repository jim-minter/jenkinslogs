package detectors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

const (
	ansiblePlaybookStartPrefix = "PLAYBOOK: "
	ansiblePlaybookEndPrefix   = "PLAY RECAP "
)

var ansiblePlaybookStartRx = regexp.MustCompile(`^PLAYBOOK: (.*) \*+$`)
var ansiblePlaybookEndRx = regexp.MustCompile(`^PLAY RECAP \*+$`)

func AnsiblePlaybookFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			if !strings.HasPrefix(s[c1], ansiblePlaybookStartPrefix) {
				continue
			}
			s1 = ansiblePlaybookStartRx.FindStringSubmatch(s[c1])
			if s1 != nil {
				c1--
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			if !strings.HasPrefix(s[c2], ansiblePlaybookEndPrefix) {
				continue
			}
			s2 = ansiblePlaybookEndRx.FindStringSubmatch(s[c2])
			if s2 != nil {
				c2++
				break
			}
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, fmt.Sprintf("ansible playbook: %s", s1[1]), false))
		}

		i = c2
	}

	return groups
}
