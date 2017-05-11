package detectors

import (
	"regexp"
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

const (
	stageStartPrefix = "########## STARTING STAGE: "
	stageEndPrefix   = "########## FINISHED STAGE: "
)

var stageStartRx = regexp.MustCompile(`^########## STARTING STAGE: (.*) ##########$`)
var stageEndRx = regexp.MustCompile(`^########## FINISHED STAGE: (SUCCESS|FAILURE): (.*) ##########$`)

func StageFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			if !strings.HasPrefix(s[c1], stageStartPrefix) {
				continue
			}
			s1 = stageStartRx.FindStringSubmatch(s[c1])
			if s1 != nil {
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			if !strings.HasPrefix(s[c2], stageEndPrefix) {
				continue
			}
			s2 = stageEndRx.FindStringSubmatch(s[c2])
			if s2 != nil {
				c2++
				break
			}
		}

		if c1 > i {
			groups = append(groups, types.NewGroup(i, c1, nil, g.CollapsedText, g.Expanded))
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, "stage: "+strings.ToLower(s1[1]), s2[1] != "SUCCESS"))
		}

		i = c2
	}

	return groups
}
