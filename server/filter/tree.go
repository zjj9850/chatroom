package filter

import "fmt"

var UTF8_TRANS []uint8 = []uint8{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5,
}

const SkipList = " .\t\r\n~!@#$%^&*()_+-=【】、[]{}|;':\"，。、《》？αβγδεζηθικλμνξοπρστυφχψωΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ。，、；：？！…—·ˉ¨‘’“”々～‖∶＂＇｀｜〃〔〕〈〉《》「」『』．〖〗【】（）［］｛｝ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩⅪⅫ⒈⒉⒊⒋⒌⒍⒎⒏⒐⒑⒒⒓⒔⒕⒖⒗⒘⒙⒚⒛㈠㈡㈢㈣㈤㈥㈦㈧㈨㈩①②③④⑤⑥⑦⑧⑨⑩⑴⑵⑶⑷⑸⑹⑺⑻⑼⑽⑾⑿⒀⒁⒂⒃⒄⒅⒆⒇≈≡≠＝≤≥＜＞≮≯∷±＋－×÷／∫∮∝∞∧∨∑∏∪∩∈∵∴⊥∥∠⌒⊙≌∽√§№☆★○●◎◇◆□℃‰€■△▲※→←↑↓〓¤°＃＆＠＼︿＿￣―♂♀┌┍┎┐┑┒┓─┄┈├┝┞┟┠┡┢┣│┆┊┬┭┮┯┰┱┲┳┼┽┾┿╀╁╂╃└┕┖┗┘┙┚┛━┅┉┤┥┦┧┨┩┪┫┃┇┋┴┵┶┷┸┹┺┻╋╊╉╈╇╆╅╄"

type FilterTree struct {
	SkipSet   map[string]struct{}
	EmptyRoot *FilterNode
}

func QueryString(text string, callback func(string, int) int) {
	for i := 0; i < len(text); i++ {
		start := text[i]
		chatLen := UTF8_TRANS[int(start)] + 1
		if int(chatLen) <= len(text)-i {
			charStr := text[i : i+int(chatLen)]
			if callback(string(charStr), i) < 0 {
				break
			}
			i += int(chatLen) - 1
		}
	}
}

func NewFilterTree() *FilterTree {
	tree := &FilterTree{}
	tree.SkipSet = make(map[string]struct{})
	tree.EmptyRoot = NewFilterNode("")
	QueryString(SkipList, func(s string, i int) int {
		tree.SkipSet[s] = struct{}{}
		return 0
	})
	return tree
}

func output_node(node *FilterNode) {
	fmt.Printf("%s ", node.CharActer)
	for _, child := range node.CharMap {
		output_node(child)
	}
}

func (self *FilterTree) OutputNodes() {
	fmt.Println("Root:")
	for _, child := range self.EmptyRoot.CharMap {
		output_node(child)
		fmt.Println("")
	}
}

func (self *FilterTree) InsertNode(text string) {
	if text == "" {
		return
	}

	var parentNode *FilterNode = self.EmptyRoot
	var childNode *FilterNode = nil

	QueryString(text, func(s string, i int) int {
		childNode = parentNode.FindChild(s)
		if childNode == nil {
			childNode = parentNode.InsertChild(s)
		}
		parentNode = childNode
		return 0
	})

	if childNode != nil {
		childNode.IsEnd = true
	}
}

func (self *FilterTree) Replace(text string) string {
	var findNode *FilterNode = self.EmptyRoot

	var check_list []string = make([]string, 0, len(text))

	QueryString(text, func(s string, i int) int {
		check_list = append(check_list, s)
		return 0
	})

	dirtyPos := make(map[int]struct{})
	skipPos := make(map[int]struct{})

	var start int = 0

	for _, ch := range check_list {
		if _, is_skip := self.SkipSet[ch]; is_skip {
			start += 1
			continue
		}

		var index int = start
		var find *FilterNode = findNode.FindChild(ch)
		for find != nil {
			if find.IsEnd {
				for j := start; j < index+1; j++ {
					dirtyPos[j] = struct{}{}
				}
			}
			if index+1 == len(check_list) {
				break
			}
			index += 1
			check_str := check_list[index]
			if _, is_skip := self.SkipSet[check_str]; !is_skip {
				find = find.FindChild(check_str)
			} else {
				skipPos[index] = struct{}{}
			}
		}
		start += 1
	}

	var newStr string = ""

	for i := 0; i < len(check_list); i++ {
		if _, dirty := dirtyPos[i]; dirty {
			if _, skip := skipPos[i]; !skip {
				newStr += "*"
			}
		} else {
			newStr += check_list[i]
		}
	}

	return newStr
}

func (self *FilterTree) Find(text string) bool {
	var findNode *FilterNode = self.EmptyRoot

	var check_list []string = make([]string, 0, len(text))

	QueryString(text, func(s string, i int) int {
		check_list = append(check_list, s)
		return 0
	})

	var start int = 0

	for _, ch := range check_list {
		if _, is_skip := self.SkipSet[ch]; is_skip {
			start += 1
			continue
		}

		var index int = start
		var foundNode *FilterNode = nil
		var find *FilterNode = findNode.FindChild(ch)
		for find != nil {
			if find.IsEnd {
				foundNode = find
			}

			if index+1 == len(check_list) {
				break
			}

			index += 1
			check_str := check_list[index]
			if _, is_skip := self.SkipSet[check_str]; !is_skip {
				find = find.FindChild(check_str)
			}
		}

		if foundNode != nil {
			return true
		}

		start += 1
	}

	return false
}
