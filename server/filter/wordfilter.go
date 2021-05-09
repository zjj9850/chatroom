package filter

type WordFilter struct {
	WordTree *FilterTree
}

func NewWordFilter() *WordFilter {
	filter := &WordFilter{}
	filter.WordTree = NewFilterTree()
	return filter
}

func (self *WordFilter) InsertKeyWords(text string) {
	self.WordTree.InsertNode(text)
}

func (self *WordFilter) ContainsAny(text string) bool {
	return self.WordTree.Find(text)
}

func (self *WordFilter) Replace(text string) string {
	return self.WordTree.Replace(text)
}

func (self *WordFilter) PrintTree() {
	self.WordTree.OutputNodes()
}
