package detectors

import "github.com/jim-minter/jenkinslogs/pkg/types"

var Detectors = []func(s []string, g *types.Group, i int, j int) []*types.Group{
	StageFind,
	AnsiblePlaybookFind,
	AnsiblePlaytaskFind,
	OSGoBuildFind,
	OSGoVetFind,
	OSEndToEndFind,
	OSIntegrationFind,
	GoTestFind,
	GinkgoFind,
	GinkgoTestFind,
	GinkgoAdditionalLoggingFind,
}
