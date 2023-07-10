package ast

import (
	"fmt"
	"static-analysis/util"
	"testing"
)

func TestEquality(t *testing.T) {
	cmp := func(lhs, rhs any) bool { return lhs == rhs }
	{
		lhs := Cons(Sexpr{"Cons"}, Sexpr{"Cell"})
		rhs := Cons(Sexpr{"Cons"}, Sexpr{"Cell"})
		if !Equals(lhs, rhs, cmp) {
			t.Errorf("Equality test failed:\n%s\n%s", lhs.Print(),
				rhs.Print())
		}
	}

	{
		lhs := S(1, 2, 3)
		rhs := S(1, 2, 3)
		if !Equals(lhs, rhs, cmp) {
			t.Errorf("Equality test failed:\n%s\n%s", lhs.Print(),
				rhs.Print())
		}
	}

	{
		lhs := S(1, S(2, "a"), 3)
		rhs := S(1, S(2, "a"), 3)
		if !Equals(lhs, rhs, cmp) {
			t.Errorf("Equality test failed:\n%s\n%s", lhs.Print(),
				rhs.Print())
		}
	}

	{
		lhs := S(1, nil)
		rhs := S(nil, 1)
		if Equals(lhs, rhs, cmp) {
			t.Errorf("Equality test failed:\n%s\n%s", lhs.Print(),
				rhs.Print())
		}
	}
}

func TestPrintDotted(t *testing.T) {
	{
		var l Sexpr = Sexpr{nil}
		expected := `nil`
		if l.PrintDotted() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.PrintDotted(), expected))
		}
	}
	{
		l := Cons(Sexpr{1}, Sexpr{2})
		expected := `(1 . 2)`
		if l.PrintDotted() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.PrintDotted(), expected))
		}
	}
	{
		l := S(1, 2, 3, 4, 5)
		expected := `(1 . (2 . (3 . (4 . (5 . nil)))))`
		if l.PrintDotted() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.PrintDotted(), expected))
		}
	}
}

func TestPrint(t *testing.T) {
	{
		l := S()
		expected := `()`
		if l.Print() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.Print(), expected))
		}
	}
	{
		l := S(1, 2, 3, 4, 5)
		expected := `(1 2 3 4 5 )`
		if l.Print() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.Print(), expected))
		}
	}
	{
		l := S("A", S("B", "C"), "D", S("E", S("F", "G")))
		expected := `(A (B C )D (E (F G )))`
		if l.Print() != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(l.Print(), expected))
		}
	}
}

func TestMinify(t *testing.T) {
	{
		l := S()
		expected := `()`
		got := Minify(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S(1, 2, 3, 4, 5)
		expected := `(1 2 3 4 5)`
		got := Minify(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("A", S("B", "C"), "D", S("E", S("F", "G")))
		expected := `(A(B C)D(E(F G)))`
		got := Minify(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
	{
		l := S("Long", "words", S("right", "here"))
		expected := `(Long words(right here))`
		got := Minify(l.Print())
		if got != expected {
			t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
		}
	}
}

func TestPrettify(t *testing.T) {
	// {
	// 	l := S()
	// 	expected := `()`
	// 	got := Indent(l.Print(), 4, 80)
	// 	if got != expected {
	// 		t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
	// 	}
	// }
	// {
	// 	l := Sexpr{"Hello"}
	// 	expected := `Hello`
	// 	got := Indent(l.Print(), 4, 80)
	// 	if got != expected {
	// 		t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
	// 	}
	// }
	// {
	// 	l := S(1, 2, 3, 4, 5)
	// 	expected := `(1 2 3 4 5)`
	// 	got := Indent(l.Print(), 4, 80)
	// 	if got != expected {
	// 		t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
	// 	}
	// }
	// {
	// 	l := S("A", S("B", "C"), "D", S("E", S("F", "G")))
	// 	expected := `(A (B C) D (E (F G)))`
	// 	got := Indent(l.Print(), 4, 80)
	// 	if got != expected {
	// 		t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
	// 	}
	// }
	// {
	// 	l := S("A", S(S("Sub", S("expr")), "C"), "D", S("E", S("F", "G")))
	// 	expected := `(A ((Sub (expr)) C) D (E (F G)))`
	// 	got := Indent(l.Print(), 4, 80)
	// 	if got != expected {
	// 		t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
	// 	}
	// }
	// {
	// 	l := S("Long", "words", S("right", "here"))
	// 	expected := `(Long words (right here))`
	// 	got := Indent(l.Print(), 4, 80)
	// 	if got != expected {
	// 		t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
	// 	}
	// }
	// {
	// 	l := S("this", "list", "contains", "very", "big", "amount", "of", "words", "as", "bad", "example", "of", "very", "long", "list")
	// 	expected := `(this list contains very big amount of words as bad example of very long list)`
	// 	got := Indent(l.Print(), 4, 80)
	// 	if got != expected {
	// 		t.Errorf("Got left, expected right\n%s\n", util.ConcatVertically(got, expected))
	// 	}
	// }
	{
		s := `(.TypeDeclaration 
	(.Attributes (.Attribute Route api/[controller] )(.Attribute ApiController ))
	(.Class TodoItemsController (.ClassBase ControllerBase )
		(.ClassMemberDeclarations 
			(.ClassMemberDeclaration (.Attributes (.Attribute HttpGet ))(.AllMemberModifiers Public Async )
				(.TypedMemberDeclaration 
					(.Type Task 
						(.TypeArgumentList 
							(.Type ActionResult 
								(.TypeArgumentList 
									(.Type IEnumerable (.TypeArgumentList (.Type TodoItem )))))))
					(.MethodDeclaration (.MethodMemberName GetTodoItems )(.FormalParameterList )(.MethodBody ))))
			(.ClassMemberDeclaration (.Attributes (.Attribute HttpGet {id} ))(.AllMemberModifiers Public Async )
				(.TypedMemberDeclaration 
					(.Type Task 
						(.TypeArgumentList (.Type ActionResult (.TypeArgumentList (.Type TodoItem )))))
					(.MethodDeclaration (.MethodMemberName GetTodoItem )
						(.FormalParameterList (.ArgDeclaration (.Type long )id ))(.MethodBody ))))
			(.ClassMemberDeclaration (.Attributes (.Attribute HttpPut {id} ))(.AllMemberModifiers Public Async )
				(.TypedMemberDeclaration (.Type Task (.TypeArgumentList (.Type IActionResult )))
					(.MethodDeclaration (.MethodMemberName PutTodoItem )
						(.FormalParameterList (.ArgDeclaration (.Type long )id )
							(.ArgDeclaration (.Type TodoItem )todoItem ))(.MethodBody ))))
			(.ClassMemberDeclaration (.Attributes (.Attribute HttpPost ))(.AllMemberModifiers Public Async )
				(.TypedMemberDeclaration 
					(.Type Task 
						(.TypeArgumentList (.Type ActionResult (.TypeArgumentList (.Type TodoItem )))))
					(.MethodDeclaration (.MethodMemberName PostTodoItem )
						(.FormalParameterList (.ArgDeclaration (.Type TodoItem )todoItem ))(.MethodBody ))))
			(.ClassMemberDeclaration (.Attributes (.Attribute HttpDelete {id} ))(.AllMemberModifiers Public Async )
				(.TypedMemberDeclaration (.Type Task (.TypeArgumentList (.Type IActionResult )))
					(.MethodDeclaration (.MethodMemberName DeleteTodoItem )
						(.FormalParameterList (.ArgDeclaration (.Type long )id ))(.MethodBody )))))))
						`
		s = `	(Source
		(FunctionDecl (main)
			(Signature (ID[]))
			(Block 
				(ConstDecl (ID[] (x)) (Expr[] (Expr (8))))
				(Expr (+ (* (x) (8)) (3)))
				(Expr (+ (x) (/ (3) (4))))
				(Expr (Call (f) 
					(Expr[] (Expr (x)) (Expr (Get (x) (y))))))
				(Assign
					(Expr[] (Expr (x)) (Expr (Get (x) (y))))
					(Expr[] (Expr (Get (x) (y))) (Expr (x))))
				
	))`
		// fmt.Println(Prettify(s))
		fmt.Println(Indent(s, 4, 80))
	}
}
