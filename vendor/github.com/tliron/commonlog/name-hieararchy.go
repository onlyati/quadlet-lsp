package commonlog

//
// NameHierarchy
//

// Convenience type for implementing maximum level per name in backends.
// Supports level inheritance.
type NameHierarchy struct {
	root *nameHierarchyNode
}

func NewNameHierarchy() *NameHierarchy {
	return &NameHierarchy{
		root: newMaxLevelHierarchyNode(),
	}
}

func (self *NameHierarchy) AllowLevel(level Level, name ...string) bool {
	return level <= self.GetMaxLevel(name...)
}

func (self *NameHierarchy) GetMaxLevel(name ...string) Level {
	node := self.root
	for _, segment := range name {
		if child, ok := node.children[segment]; ok {
			node = child
		} else {
			break
		}
	}
	return node.maxLevel
}

func (self *NameHierarchy) SetMaxLevel(level Level, name ...string) {
	node := self.root
	for _, segment := range name {
		if child, ok := node.children[segment]; ok {
			node = child
		} else {
			child = newMaxLevelHierarchyNode()
			node.children[segment] = child
			node = child
		}
	}
	node.maxLevel = level
}

//
// nameHierarchyNode
//

type nameHierarchyNode struct {
	maxLevel Level
	children map[string]*nameHierarchyNode
}

func newMaxLevelHierarchyNode() *nameHierarchyNode {
	return &nameHierarchyNode{
		maxLevel: None,
		children: make(map[string]*nameHierarchyNode),
	}
}
