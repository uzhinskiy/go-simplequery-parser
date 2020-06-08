# go-simplequery-parser

Parse a condition query into a nested struct.

## Example

```golang
root, _ := Parse(`foo < 10 AND (bar = "x" OR NOT baz = "y")`)
```

returns a root Node that has the following nested structure:

```golang
AND {
  LT {
    ID {
      "foo",
    },
    VAL {
      "10",
    },
  },
  OR {
    EQ {
      ID {
        "bar",
      },
      VAL {
        "x",
      },
    },
    NOT {
      EQ {
        ID {
          "baz",
        },
        VAL {
          "y",
        },
      },
    },
  },
}
```

This package does not implement a matcher that checks data against it, but this can be easily achieved with functions like

```golang
func Match(node Node, data <your data>) bool {
    switch n := node.(type) {
    case AND:
        return Match(n.Node1, data) && Match(n.Node2, data)
    case OR:
        return Match(n.Node1, data) || Match(n.Node2, data)
    case NOT:
        return !Match(n.Node, data)
    ...
    case EQ:
        return n.Note1.NodeValue == n.Note2.NodeValue
    case LT:
        ....
    ....
    }
}
```

Alternatively this package may be useful to generate another search query with another syntax from it *(Atm I use this package to generate a MongoDB search query based on the above simpler query syntax)*.