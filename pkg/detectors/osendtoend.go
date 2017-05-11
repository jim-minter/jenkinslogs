package detectors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

const (
	osEndToEndStartPrefix      = "Running "
	osEndToEndSuccessEndPrefix = "SUCCESS after "
	osEndToEndFailureEndPrefix = "FAILURE after "
)

var osEndToEndStartRx = regexp.MustCompile(`^Running [^ ]+\.sh:\d+: executing (.+)\.\.\.$`)
var osEndToEndEndRx = regexp.MustCompile(`^(SUCCESS|FAILURE) after \d+\.\d+s`)

func OSEndToEndFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			if !strings.HasPrefix(s[c1], osEndToEndStartPrefix) {
				continue
			}
			s1 = osEndToEndStartRx.FindStringSubmatch(s[c1])
			if s1 != nil {
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			if !strings.HasPrefix(s[c2], osEndToEndSuccessEndPrefix) && !strings.HasPrefix(s[c2], osEndToEndFailureEndPrefix) {
				continue
			}
			s2 = osEndToEndEndRx.FindStringSubmatch(s[c2])
			if s2 != nil {
				c2++
				break
			}
		}

		if s1 != nil && s2 != nil {
			groups = append(groups, types.NewGroup(c1, c2, nil, fmt.Sprintf("os end to end: %s", s1[1]), s2[1] != "SUCCESS"))
		}

		i = c2
	}

	return groups
}
