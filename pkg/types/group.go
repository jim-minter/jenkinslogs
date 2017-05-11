package types

type Group struct {
	ID            int
	I, J          int
	Children      []*Group
	CollapsedText string
	Expanded      bool
}

var id int

func NewGroup(i, j int, children []*Group, collapsedText string, expanded bool) *Group {
	id++
	return &Group{ID: id, I: i, J: j, Children: children, CollapsedText: collapsedText, Expanded: expanded}
}

func (g *Group) Walk(s []string, enterGroup func(*Group), exitGroup func(*Group), f func([]string, *Group, int, int) []*Group) {
	children := []*Group{}

	if enterGroup != nil {
		enterGroup(g)
	}

	c := g.I
	for _, child := range g.Children {
		if f != nil && c < child.I {
			children = append(children, f(s, g, c, child.I)...)
		}
		child.Walk(s, enterGroup, exitGroup, f)
		children = append(children, child)
		c = child.J
	}
	if f != nil && c < g.J {
		children = append(children, f(s, g, c, g.J)...)
	}

	if exitGroup != nil {
		exitGroup(g)
	}

	g.Children = children
}
