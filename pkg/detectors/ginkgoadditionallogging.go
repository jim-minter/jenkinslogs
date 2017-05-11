package detectors

import (
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

const (
	ginkgoAdditionalLoggingStartPrefix = "STEP: Collecting events from namespace"
	ginkgoAdditionalLoggingEnd         = "STEP: Dumping a list of prepulled images on each node"
)

func GinkgoAdditionalLoggingFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			if strings.HasPrefix(s[c1], ginkgoAdditionalLoggingStartPrefix) {
				s1 = []string{}
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			if s[c2] == ginkgoAdditionalLoggingEnd {
				s2 = []string{}
				c2++
				break
			}
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, "ginkgo additional logging", false))
		}

		i = c2
	}

	return groups
}
