package simplequery

import (
	"fmt"
	"strings"
)

func ExampleParse() {
	root, err := Parse(`foo < 10 AND (bar = "x" OR NOT baz = "y")`)
	if err != nil {
		panic(err)
	}
	prettyPrint(root)
	// Output:
	// AND {
	//   LT {
	//     ID {
	//       "foo",
	//     },
	//     VAL {
	//       "10",
	//     },
	//   },
	//   OR {
	//     EQ {
	//       ID {
	//         "bar",
	//       },
	//       VAL {
	//         "x",
	//       },
	//     },
	//     NOT {
	//       EQ {
	//         ID {
	//           "baz",
	//         },
	//         VAL {
	//           "y",
	//         },
	//       },
	//     },
	//   },
	// }
}

func prettyPrint(root Node) {
	str := root.String()
	var indent = 0
	for _, char := range str {
		if char == '{' {
			fmt.Println(" {")
			indent++
			fmt.Print(strings.Repeat("  ", indent))
		} else if char == '}' {
			fmt.Println(",")
			indent--
			fmt.Print(strings.Repeat("  ", indent), "}")
		} else if char == ',' {
			fmt.Println(",")
			fmt.Print(strings.Repeat("  ", indent))
		} else {
			fmt.Print(string(char))
		}
	}
}
