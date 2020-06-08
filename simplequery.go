package simplequery

func Parse(condition string) (Node, error) {
	tokens, err := tokenize(condition)
	if err != nil {
		return nil, err
	}
	return parse(tokens)
}
