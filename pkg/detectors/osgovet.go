package detectors

import (
	"regexp"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

var osGoVetStartRx = regexp.MustCompile(`^[^ ]+\.go:\d+: `)
var osGoVetEndRx = regexp.MustCompile(`^(FAILURE|SUCCESS): go vet (?:failed|succeeded)!$`)

func OSGoVetFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c2 = i; c2 < j; c2++ {
			if s[c2] != "SUCCESS: go vet succeeded!" && s[c2] != "FAILURE: go vet failed!" {
				continue
			}
			s2 = osGoVetEndRx.FindStringSubmatch(s[c2])
			if s2 != nil {
				c2++
				break
			}
		}

		if s2 != nil {
			for c1 = c2 - 1; c1 > i; c1-- {
				s1 = osGoVetEndRx.FindStringSubmatch(s[c1-1])
				if s1 == nil {
					break
				}
			}

			groups = append(groups, types.NewGroup(c1, c2, nil, "os go vet", s2[1] != "SUCCESS"))
		}

		i = c2
	}

	return groups
}
