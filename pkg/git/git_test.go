package git

import (
	"testing"
)

const (
	testFilePath = "../../test_www/gitfolder/index"
)

var g *Git

func TestNew(t *testing.T) {
	var err error
	g, err = New(testFilePath)
	if err != nil {
		t.Error(err)
	}
}

func TestParseIndex(t *testing.T) {
	if g.Index.Header.Signature != "DIRC" {
		t.Error("sinagure is invalid")
	}
}

func TestCheckNullBytes(t *testing.T) {
	var (
		isNulled    = []byte{0, 0, 0, 0}
		isNonNulled = []byte{0, 1, 0, 0}
	)

	nulledTest := g.checkNullBytes(isNulled)
	if !nulledTest {
		t.Error("nulled bytes check fail")
	}

	nonNulledTest := g.checkNullBytes(isNonNulled)
	if nonNulledTest {
		t.Error("nonNulled bytes check fail")
	}
}

func TestReadBytes(t *testing.T) {
	if g.pos != 209 {
		t.Error("cursor pos is not 0")
	}
	g.readBytes(5)
	if g.pos != 214 {
		t.Error("byte read fail")
	}

}
