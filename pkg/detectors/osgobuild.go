package detectors

import (
	"regexp"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

var osGoBuildStartRx = regexp.MustCompile(`^(?:\++ )?Building go targets for `)
var osGoBuildEndRx = regexp.MustCompile(`^(?:# )|(?:[^ ]+\.go:\d+: )`)

func OSGoBuildFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			s1 = osGoBuildStartRx.FindStringSubmatch(s[c1])
			if s1 != nil {
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			s2 = osGoBuildEndRx.FindStringSubmatch(s[c2])
			if s2 == nil {
				s2 = []string{s[c2]}
				break
			}
			s2 = nil
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, "os go build", c2-c1 > 1))
		}

		i = c2
	}

	return groups
}
