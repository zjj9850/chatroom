package filter

type FilterNode struct {
	CharActer string
	CharMap   map[string]*FilterNode
	IsEnd     bool
}

func NewFilterNode(char string) *FilterNode {
	node := &FilterNode{
		CharActer: char,
		IsEnd:     false,
	}
	node.CharMap = make(map[string]*FilterNode)
	return node
}

func (self *FilterNode) FindChild(char string) *FilterNode {
	if c, e := self.CharMap[char]; e {
		return c
	}
	return nil
}

func (self *FilterNode) InsertChild(char string) *FilterNode {
	node := self.FindChild(char)
	if node == nil {
		node = NewFilterNode(char)
		self.CharMap[char] = node
	}
	return node
}
