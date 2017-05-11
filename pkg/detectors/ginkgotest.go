package detectors

import (
	"regexp"
	"strings"

	"github.com/jim-minter/jenkinslogs/pkg/types"
)

const (
	ginkgoTestDivider       = "------------------------------"
	ginkgoTestFailurePrefix = "• Failure"
)

var ginkgoTestSkippedRx = regexp.MustCompile(`^S+(?:$| \[SKIPPING\])`)
var ginkgoTestNameEndRx = regexp.MustCompile(`\.go:\d+$`)

func GinkgoTestFind(s []string, g *types.Group, i int, j int) []*types.Group {
	groups := []*types.Group{}

	for i < j {
		var s1, s2 []string
		var c1, c2 int

		for c1 = i; c1 < j; c1++ {
			if s[c1] == ginkgoTestDivider {
				s1 = []string{}
				break
			}
		}

		for c2 = c1 + 1; c2 < j; c2++ {
			if s[c2] == ginkgoTestDivider {
				s2 = []string{}
				break
			}
		}

		if s1 != nil && s2 != nil {
			testName := []string{}
			readingTestName := true
			failed := false
			skipped := false
			for c := c1 + 1; c < c2; c++ {
				if readingTestName {
					readingTestName = !ginkgoTestNameEndRx.MatchString(s[c])
				}
				if readingTestName {
					testName = append(testName, s[c])
				}
				if !failed {
					failed = strings.HasPrefix(s[c], ginkgoTestFailurePrefix)
				}
				if !skipped {
					skipped = skipped || ginkgoTestSkippedRx.MatchString(s[c])
				}
			}

			var collapsedText string
			if skipped {
				collapsedText = "skipped ginkgo test(s)"
			} else {
				collapsedText = "ginkgo test: " + strings.Join(testName, " / ")
			}

			groups = append(groups, types.NewGroup(c1, c2, nil, collapsedText, failed))
		}

		i = c2
	}

	return groups
}
