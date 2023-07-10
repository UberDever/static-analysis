// Following implementation of s-expressions assumes that they are trees (acyclic graphs)
package ast

import (
	"fmt"
	"sort"
	"static-analysis/util"
	"strings"
	"unicode"

	Text "github.com/linkdotnet/golang-stringbuilder"
)

type Sexpr struct {
	any
}

type cons struct {
	lhs, rhs Sexpr
}

func (v Sexpr) IsAtom() bool {
	switch v.any.(type) {
	case nil:
		return true
	case bool:
		return true
	case uintptr:
		return true
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case float32:
		return true
	case float64:
		return true
	case complex64:
		return true
	case complex128:
		return true
	case string:
		return true
	default:
		return false
	}
}

func (v Sexpr) IsNil() bool {
	return v.any == nil
}

func (v Sexpr) PrintDotted() string {
	s := strings.Builder{}
	if v.IsAtom() {
		switch val := v.any.(type) {
		case nil:
			s.WriteString("nil")
		default:
			s.WriteString(fmt.Sprint(val))
		}
	} else {
		s.WriteByte('(')
		s.WriteString(Car(v).PrintDotted())
		s.WriteString(" . ")
		s.WriteString(Cdr(v).PrintDotted())
		s.WriteByte(')')
	}
	return s.String()
}

func (v Sexpr) Print() string {
	if v.IsNil() {
		return "()"
	}

	s := strings.Builder{}
	onEnter := func(node Sexpr) {
		if node.IsAtom() {
			s.WriteString(fmt.Sprint(node.any))
		} else {
			s.WriteByte('(')
		}
	}
	onExit := func(node Sexpr) {
		if node.IsAtom() {
			s.WriteByte(' ')
		} else {
			s.WriteByte(')')
		}
	}
	TraversePreorder(v, onEnter, onExit)
	return s.String()
}

func S(list ...any) Sexpr {
	if len(list) == 0 {
		return Sexpr{nil}
	}
	expr := Sexpr{}
	switch list[0].(type) {
	case Sexpr:
		expr = list[0].(Sexpr)
	default:
		expr.any = list[0]
	}
	return Cons(expr, S(list[1:]...))
}

func Cons(lhs any, rhs any) Sexpr {
	var l, r Sexpr
	switch lhs := lhs.(type) {
	case Sexpr:
		l = lhs
	default:
		l = Sexpr{lhs}
	}
	switch rhs := rhs.(type) {
	case Sexpr:
		r = rhs
	default:
		r = Sexpr{rhs}
	}
	return Sexpr{cons{lhs: l, rhs: r}}
}

func Car(v Sexpr) Sexpr {
	switch v.any.(type) {
	case cons:
		cons := v.any.(cons)
		return cons.lhs
	default:
		return Sexpr{nil}
	}
}

func Cdr(v Sexpr) Sexpr {
	switch v.any.(type) {
	case cons:
		cons := v.any.(cons)
		return cons.rhs
	default:
		return Sexpr{nil}
	}
}

func Equals(lhs Sexpr, rhs Sexpr, cmp func(any, any) bool) bool {
	if lhs.IsAtom() || rhs.IsAtom() {
		return cmp(lhs.any, rhs.any)
	}

	return Equals(Car(lhs), Car(rhs), cmp) &&
		Equals(Cdr(lhs), Cdr(rhs), cmp)
}

func Minify(sexpr string) string {
	if !strings.ContainsRune(sexpr, '(') {
		return strings.TrimSpace(sexpr)
	}

	s := Text.StringBuilder{}
	wasSpace := false
	for _, c := range sexpr {
		if c == ')' || c == '(' {
			s.AppendRune(c)
		} else if unicode.IsSpace(c) {
			wasSpace = true
			continue
		} else {
			last := s.RuneAt(s.Len() - 1)
			if last != ')' && last != '(' && wasSpace {
				s.AppendRune(' ')
			}
			s.AppendRune(c)
		}
		wasSpace = false
	}
	return s.ToString()
}

func Prettify(sexpr string) string {
	if !strings.ContainsRune(sexpr, '(') {
		return strings.TrimSpace(sexpr)
	}
	minified := Minify(sexpr)

	s := strings.Builder{}
	for i := 0; i < len(minified); i++ {
		c := minified[i]
		if c != '(' && c != ')' && c != ' ' {
			var prevClose, nextOpen bool
			prevClose = minified[i-1] == ')'
			nextOpen = minified[i+1] == '('
			if prevClose {
				s.WriteByte(' ')
			}
			s.WriteByte(c)
			if nextOpen {
				s.WriteByte(' ')
			}
		} else {
			if i > 0 && minified[i-1] == ')' && c == '(' {
				s.WriteByte(' ')
			}
			s.WriteByte(c)
		}
	}
	return s.String()
}

type indentInfo = struct{ begin, end, depth int }

func calculateIndentation(sexpr string, width, cutoff int) []indentInfo {
	rawInfo := util.NewStack[indentInfo]()
	openParens := util.NewStack[int]()

	curDepth := -1
	for i := 0; i < len(sexpr); i++ {
		if sexpr[i] == '(' {
			openParens.Push(i)
			curDepth++
		} else if sexpr[i] == ')' {
			begin := openParens.ForcePop()
			end := i
			rawInfo.Push(indentInfo{begin, end, curDepth})
			curDepth--
		}
	}

	resultInfo := make([]indentInfo, 0, 64)
	lineEnd := len(sexpr)
	for !rawInfo.IsEmpty() {
		pos := rawInfo.ForcePop()
		exprLen := (pos.end - pos.begin + 1)
		newOffset := pos.depth * width
		if lineEnd+exprLen >= cutoff && pos.depth > 0 {
			resultInfo = append(resultInfo, indentInfo{pos.begin, pos.end, pos.depth})
			lineEnd = newOffset + exprLen
		}
	}
	sort.SliceStable(resultInfo, func(i, j int) bool { return resultInfo[i].begin < resultInfo[j].begin })
	return resultInfo
}

func Indent(sexpr string, width, cutoff int) string {
	if !strings.ContainsRune(sexpr, '(') {
		return strings.TrimSpace(sexpr)
	}

	pretty := Prettify(sexpr)
	indicies := calculateIndentation(pretty, width, cutoff)
	if len(indicies) == 0 {
		return pretty
	}

	s := strings.Builder{}
	curIdx := 0
	idx := indicies[curIdx]
	for i, c := range pretty {
		if i == idx.begin {
			s.WriteByte('\n')
			newOffset := idx.depth * width
			for j := 0; j < newOffset; j++ {
				s.WriteByte(' ')
			}
			curIdx++
			if curIdx < len(indicies) {
				idx = indicies[curIdx]
			}
		}
		s.WriteRune(c)
	}
	return s.String()
}

type Action func(Sexpr)

func TraversePreorder(root Sexpr, onEnter Action, onExit Action) {
	traversePreorderRec(onEnter, onExit, root)
}

func traversePreorderRec(onEnter Action, onExit Action, cur Sexpr) {
	if cur.any == nil {
		return
	}

	onEnter(cur)
	defer onExit(cur)
	children := cur
	for c := Car(children); c.any != nil; c = Car(children) {
		children = Cdr(children)
		traversePreorderRec(onEnter, onExit, c)
	}
}

func TraversePostorder(root Sexpr, onEnter Action) {
	traversePostorderRec(onEnter, root)
}

func traversePostorderRec(onEnter Action, cur Sexpr) {
	c := Car(cur)
	if c.any == nil {
		return
	}

	traversePostorderRec(onEnter, Cdr(cur))
	traversePostorderRec(onEnter, c)
	onEnter(c)
}
