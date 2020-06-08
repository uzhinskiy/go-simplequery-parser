package simplequery

import "fmt"

func Example_tokenize() {
	tokens, err := tokenize(`foo < 10 AND (bar = "x" OR NOT baz = "y")`)
	if err != nil {
		panic(err)
	}
	for _, token := range tokens {
		fmt.Println(token)
	}
	// Output:
	// ID["foo"]{1}
	// LT["<"]{5}
	// VAL["10"]{7}
	// AND["AND"]{10}
	// LPAR["("]{14}
	// ID["bar"]{15}
	// EQ["="]{19}
	// VAL["x"]{21}
	// OR["OR"]{25}
	// NOT["NOT"]{28}
	// ID["baz"]{32}
	// EQ["="]{36}
	// VAL["y"]{38}
	// RPAR[")"]{41}
}
