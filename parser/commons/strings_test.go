package commons_test

import (
	"testing"

	"github.com/pharosnet/dalg/parser/commons"
)

func TestSnakeToCamel(t *testing.T) {
	s := "aggregateName"
	s = commons.SnakeToCamel(s)
	t.Log(s)
}
