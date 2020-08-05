package commons_test

import (
	"testing"

	"github.com/pharosnet/dalg/parser/commons"
)

func TestWordsContainsAll(t *testing.T) {
	words := []string{"abc", "ddd", "---", "123"}
	has := commons.WordsContainsAll(words, "---", "ddd")
	t.Log(has)
	t.Log(commons.WordsContainsAll(words, "---", "xx"))
}
