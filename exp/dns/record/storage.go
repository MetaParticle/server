package record

type Tree interface {
	//Walks through the tree until last string or reaching a leaf.
	Walk([]string) Tree
	
	//Returns the parent of the tree, or nil if tree is the root.
	Parent() Tree
	
	//returns the root of the tree.
	Root() Tree
	
	//Adds the child string to the tree, containing data.
	Add(string, interface{})
	
	//Removes the child from the tree, and all of its children.
	Remove(string)
	
	//Returns the data stored in the tree.
	Data() interface{}
}

type MapTree struct {
	children map[string]Tree
	parent Tree
	data interface{}
}

func (mt MapTree) Walk(addr []string) Tree {
	for _, child := range addr {
		mt = mt.children[]
	}
	return mt
}

func (mt MapTree) Parent() Tree {
	return mt.parent
}

func (mt MapTree) Root() Tree {
	for root := Tree(mt); root != nil {
		root = root.Parent()
	}
}

func (mt MapTree) Add(node string, data interface) {
	mt.children[node] = data
}

func (mt MapTree) Remove(node string) {
	delete(mt, node)
}

func (mt MapTree) Data() interface{} {
	return mt.data
}
