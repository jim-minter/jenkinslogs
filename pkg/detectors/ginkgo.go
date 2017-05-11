package detectors

import (
	"regexp"
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

const (
	ginkgoStartPrefix      = "Running Suite: "
	ginkgoSuccessEndPrefix = "SUCCESS! -- "
	ginkgoFailEndPrefix    = "FAIL! -- "
)

var ginkgoEndRx = regexp.MustCompile(`^(SUCCESS|FAIL)! -- \d+ Passed \| \d+ Failed \| \d+ Pending \| \d+ Skipped`)

func GinkgoFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			if strings.HasPrefix(s[c1], ginkgoStartPrefix) {
				s1 = []string{}
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			if !strings.HasPrefix(s[c2], ginkgoSuccessEndPrefix) && !strings.HasPrefix(s[c2], ginkgoFailEndPrefix) {
				continue
			}
			s2 = ginkgoEndRx.FindStringSubmatch(s[c2])
			if s2 != nil {
				c2++
				break
			}
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, "ginkgo suite", s2[1] != "SUCCESS"))
		}

		i = c2
	}

	return groups
}
