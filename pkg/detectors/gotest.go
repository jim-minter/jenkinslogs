package detectors

import (
	"regexp"
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

const (
	goTestStartPrefix = "=== RUN   Test"
	goTestEndPrefix1  = "--- "
	goTestEndPrefix2  = "FAIL\t"
)

var goTestStartRx = regexp.MustCompile(`^=== RUN   (Test.*)`)
var goTestEndRx1 = regexp.MustCompile(`^--- ((?:PASS)|(?:FAIL)): Test`)
var goTestEndRx2 = regexp.MustCompile(`^(FAIL)\t`)

func GoTestFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			if !strings.HasPrefix(s[c1], goTestStartPrefix) {
				continue
			}
			s1 = goTestStartRx.FindStringSubmatch(s[c1])
			if s1 != nil {
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			if !strings.HasPrefix(s[c2], goTestEndPrefix1) && !strings.HasPrefix(s[c2], goTestEndPrefix2) {
				continue
			}
			s2 = goTestEndRx1.FindStringSubmatch(s[c2])
			if s2 != nil {
				c2++
				break
			}
			s2 = goTestEndRx2.FindStringSubmatch(s[c2])
			if s2 != nil {
				break
			}
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, "go test: "+s1[1], s2[1] != "PASS"))
		}

		i = c2
	}

	return groups
}
