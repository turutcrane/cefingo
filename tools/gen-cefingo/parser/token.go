package parser

import (
	"log"
	"path/filepath"
	"regexp"
	"strings"

	strcase "github.com/stoewer/go-strcase"
	"modernc.org/xc"
)

type Token xc.Token

var noToken = Token{}

func (t Token) Name() string {
	return string(xc.Token(t).S())
}

func (t Token) String() string {
	return t.Name() + ": " + t.FilePos()
}

// BaseName _cef_app_t -> app
// TitleCase AppT
// GoName    CAppT

var reBase = regexp.MustCompile("^_?cef_(.*)_t$")
var reBaseT = regexp.MustCompile("^_?cef_(.*)$")

func (t Token) BaseName() string {

	match := reBase.FindStringSubmatch(t.Name())
	if match != nil {
		return match[1]
	}
	return t.Name()
}

func (t Token) TitleCase() string {
	match := reBaseT.FindStringSubmatch(t.Name())
	if match != nil {
		return strcase.UpperCamelCase(match[1])
	}
	return strcase.UpperCamelCase(t.Name())
}

func (t Token) GoName() string {
	return "C" + t.TitleCase()
}

func (t Token) Line() int {
	return xc.FileSet.Position(xc.Token(t).Pos()).Line
}

func relFilename(fn string) string {
	rel, err := filepath.Rel(cefdir, fn)
	if err != nil {
		log.Panicf("T55: can not take rel: %s, %v\n", fn, err)
	}
	return filepath.ToSlash(rel)
}

func (t Token) AbsFilename() string {
	xcToken := xc.Token(t)
	fn := (&xcToken).Position().Filename
	return fn
}

func (t Token) Filename() string {
	xcToken := xc.Token(t)
	fn := (&xcToken).Position().Filename
	return relFilename(fn)
}

func (t Token) FilePos() string {
	xcToken := xc.Token(t)
	fn := (&xcToken).Position().Filename
	pos := (&xcToken).Position().String()
	pos = strings.Replace(pos, fn, "", 1)
	rel := relFilename(fn)
	return rel + pos
}
