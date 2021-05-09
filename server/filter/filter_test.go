package filter_test

import (
	"chatroom/filter"
	"math/rand"
	"testing"
	"time"
)

var (
	DirtyWords = []string{
		"ar5e", "arrse", "arse", "ass", "ass-fucker",
		"asses", "assfucker", "assfukka", "asshole",
		"b17ch", "b1tch", "ballbag", "balls", "ballsack",
		"bastard", "beastial", "beastiality", "bellend", "bestial",
		"bestiality", "bi+ch", "biatch", "bitch", "bitcher", "bitchers",
		"blowjobs", "boiolas", "bollock", "bollok", "boner", "boob", "boobs",
		"booobs", "boooobs", "booooobs", "booooooobs", "breasts",
	}

	Letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randString() string {
	b := make([]rune, 8)
	for i := range b {
		b[i] = Letters[rand.Intn(len(Letters))]
	}
	return string(b)
}

func TestNewWordFilter(t *testing.T) {
	filter.NewWordFilter()
}

func TestInsertKeyWords(t *testing.T) {
	wordFilter := filter.NewWordFilter()

	for _, key := range DirtyWords {
		wordFilter.InsertKeyWords(key)
	}
}

func TestContainsAny(t *testing.T) {
	wordFilter := filter.NewWordFilter()

	for _, key := range DirtyWords {
		wordFilter.InsertKeyWords(key)
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		idx := rand.Intn(len(DirtyWords))
		words := randString() + DirtyWords[idx] + randString()
		if wordFilter.ContainsAny(words) {
			t.Log("Words:", words, " DirtyWords", DirtyWords[idx])
		}
	}
}

func TestReplace(t *testing.T) {
	wordFilter := filter.NewWordFilter()

	for _, key := range DirtyWords {
		wordFilter.InsertKeyWords(key)
	}

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		idx := rand.Intn(len(DirtyWords))
		words := randString() + DirtyWords[idx] + randString()
		t.Log("Replace ", words, " To ", wordFilter.Replace(words))
	}

}

func TestPrintTree(t *testing.T) {
	wordFilter := filter.NewWordFilter()

	for _, key := range DirtyWords {
		wordFilter.InsertKeyWords(key)
	}

	wordFilter.PrintTree()
}
