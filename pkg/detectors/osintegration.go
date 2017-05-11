package detectors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

var osIntegrationStartRx = regexp.MustCompile(`^Running (Test.+)\.\.\.$`)
var osIntegrationEndRx = regexp.MustCompile(`^(ok    |failed)  Test`)

func OSIntegrationFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			// s1 = osIntegrationStartRx.FindStringSubmatch(s[c1])
			if strings.HasPrefix(s[c1], "Running Test") && strings.HasSuffix(s[c1], "...") {
				s1 = []string{s[c1], s[c1][8 : len(s[c1])-3]}
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			s2 = osIntegrationEndRx.FindStringSubmatch(s[c2])
			if s2 != nil {
				c2++
				break
			}
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, fmt.Sprintf("os integration: %s", s1[1]), s2[1] != "ok    "))
		}

		i = c2
	}

	return groups
}
