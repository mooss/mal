package reader

import (
	"fmt"
	"strconv"

	hcore "github.com/mooss/harper/go/core"
	"github.com/mooss/harper/go/rad"
	"github.com/mooss/mal/src/core"
	"github.com/mooss/mal/src/lexer"
	"github.com/mooss/mal/src/strfun"
	"github.com/mooss/mal/src/typ"
	"github.com/mooss/mal/src/ut"
)

var Dialect = rad.Dialect[lexer.Kind, core.Value]{Rules: []rad.Rule[lexer.Kind, core.Value]{
	rad.Atom(lexer.KindInt, readInt),
	rad.Atom(lexer.KindDQS, readString),
	rad.Atom(lexer.KindSymbol, readSymbol),
	rad.Sequence(lexer.KindDelimiter, "(", ")", buildList),
	rad.Sequence(lexer.KindDelimiter, "[", "]", buildVector),
	rad.Table(lexer.KindDelimiter, "{", "}", 2, buildMap),
	rad.Prefix(lexer.KindSigil, "'", wrapList("quote")),
	rad.Prefix(lexer.KindSigil, "`", wrapList("quasiquote")),
	rad.Prefix(lexer.KindSigil, "~", wrapList("unquote")),
	rad.Prefix(lexer.KindSigil, "@", wrapList("deref")),
	rad.Prefix(lexer.KindSigil, "~@", wrapList("splice-unquote")),
}}

var Spec = ut.Must(Dialect.Compile())

func ReadStr(input string) (*rad.Reader[lexer.Kind, core.Value], error) {
	toks, err := lexer.Tokenize(input)
	if err != nil {
		return nil, err
	}
	return Spec.NewReader(toks), nil
}

func ParseOne(input string) (core.Value, error) {
	toks, err := lexer.Tokenize(input)
	if err != nil {
		return nil, err
	}
	return Spec.ParseOne(toks)
}

func ParseAll(input string) ([]core.Value, error) {
	toks, err := lexer.Tokenize(input)
	if err != nil {
		return nil, err
	}
	return Spec.ParseAll(toks)
}

//////////////
// Builders //

func buildList(v []core.Value) (core.Value, error) { return typ.List(v), nil }

func buildMap(v []core.Value) (core.Value, error) {
	assocs := make([]typ.Assoc, 0, len(v)/2)
	for i := 0; i < len(v); i += 2 {
		assocs = append(assocs, typ.Assoc{Key: v[i], Value: v[i+1]})
	}
	return typ.MapForm(assocs), nil
}

func buildVector(v []core.Value) (core.Value, error) { return typ.VectorForm(v), nil }

func readInt(tok hcore.Token[lexer.Kind]) (core.Value, error) {
	i, err := strconv.ParseInt(tok.Value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot read %q as int: %w", tok.Value, err)
	}
	return typ.Integer(i), nil
}
func readString(tok hcore.Token[lexer.Kind]) (core.Value, error) {
	s, err := strfun.MalUnquote(tok.Value)
	if err != nil {
		return nil, err
	}
	return typ.NewString(s), nil
}

func readSymbol(tok hcore.Token[lexer.Kind]) (core.Value, error) {
	switch tok.Value {
	case "nil":
		return typ.Nil{}, nil
	case "true":
		return typ.Bool(true), nil
	case "false":
		return typ.Bool(false), nil
	}
	if len(tok.Value) > 0 && tok.Value[0] == ':' {
		return typ.NewKeyword(tok.Value[1:]), nil
	}
	return typ.NewSymbol(tok.Value), nil
}

func wrapList(sym string) func(core.Value) (core.Value, error) {
	return func(n core.Value) (core.Value, error) {
		return typ.NewList(typ.NewSymbol(sym), n), nil
	}
}
