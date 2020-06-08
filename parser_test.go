package simplequery

import (
	"fmt"
)

func lexAndParse(condition string) {
	tokens, err := tokenize(condition)
	if err != nil {
		panic(err)
	}
	node, err := parse(tokens)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(node.String())
	}
}

func Example_parse() {
	lexAndParse("X = 5")
	lexAndParse("X = 5 AND Y = 6")
	lexAndParse("(X = 5 AND Y = 6)")
	lexAndParse("(X = 5 AND Y = 6) OR (X = 10 AND Y = 11)")
	lexAndParse("X < 5")
	lexAndParse("X <= 5")
	lexAndParse("NOT X = 5")
	lexAndParse("NOT (X = 10 AND Y = 11)")
	// Output:
	// EQ{ID{"X"},VAL{"5"}}
	// AND{EQ{ID{"X"},VAL{"5"}},EQ{ID{"Y"},VAL{"6"}}}
	// AND{EQ{ID{"X"},VAL{"5"}},EQ{ID{"Y"},VAL{"6"}}}
	// OR{AND{EQ{ID{"X"},VAL{"5"}},EQ{ID{"Y"},VAL{"6"}}},AND{EQ{ID{"X"},VAL{"10"}},EQ{ID{"Y"},VAL{"11"}}}}
	// LT{ID{"X"},VAL{"5"}}
	// LTE{ID{"X"},VAL{"5"}}
	// NOT{EQ{ID{"X"},VAL{"5"}}}
	// NOT{AND{EQ{ID{"X"},VAL{"10"}},EQ{ID{"Y"},VAL{"11"}}}}
}
