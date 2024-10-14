package simplequery

import (
	"errors"
	"fmt"
)

type Node interface {
	String() string
	Value() string
	Children() (Node, Node)
}

type (
	oneSubNode  struct{ node Node }
	twoSubNodes struct{ node1, node2 Node }
	valueNode   struct{ nodeValue string }
)

func (_ oneSubNode) Value() string  { return "" }
func (_ twoSubNodes) Value() string { return "" }
func (n valueNode) Value() string   { return n.nodeValue }

func (n oneSubNode) Children() (Node, Node)  { return n.node, nil }
func (n twoSubNodes) Children() (Node, Node) { return n.node1, n.node2 }
func (_ valueNode) Children() (Node, Node)   { return nil, nil }

type (
	AND struct{ twoSubNodes }
	OR  struct{ twoSubNodes }
	NOT struct{ oneSubNode }
	LT  struct{ twoSubNodes }
	LTE struct{ twoSubNodes }
	GT  struct{ twoSubNodes }
	GTE struct{ twoSubNodes }
	EQ  struct{ twoSubNodes }
	SET  struct{ twoSubNodes }
	NE  struct{ twoSubNodes }
	ID  struct{ valueNode }
	VAL struct{ valueNode }
)

func (n AND) String() string { return fmt.Sprintf("AND{%s,%s}", n.node1, n.node2) }
func (n OR) String() string  { return fmt.Sprintf("OR{%s,%s}", n.node1, n.node2) }
func (n NOT) String() string { return fmt.Sprintf("NOT{%s}", n.node) }
func (n LT) String() string  { return fmt.Sprintf("LT{%s,%s}", n.node1, n.node2) }
func (n LTE) String() string { return fmt.Sprintf("LTE{%s,%s}", n.node1, n.node2) }
func (n GT) String() string  { return fmt.Sprintf("GT{%s,%s}", n.node1, n.node2) }
func (n GTE) String() string { return fmt.Sprintf("GTE{%s,%s}", n.node1, n.node2) }
func (n EQ) String() string  { return fmt.Sprintf("EQ{%s,%s}", n.node1, n.node2) }
func (n SET) String() string  { return fmt.Sprintf("SET{%s,%s}", n.node1, n.node2) }
func (n NE) String() string  { return fmt.Sprintf("NE{%s,%s}", n.node1, n.node2) }
func (n ID) String() string  { return fmt.Sprintf("ID{%q}", n.nodeValue) }
func (n VAL) String() string { return fmt.Sprintf("VAL{%q}", n.nodeValue) }

func parse(tokens []token) (Node, error) {
	res := make([]interface{}, len(tokens))
	for i, token := range tokens {
		res[i] = token
	}
	// identify subexpressions with parentheses
	for {
		lPos := findToken(tokens, tokenTypeLPAR)
		if lPos < 0 {
			break
		}
		if lPos >= len(tokens)-1 {
			return nil, errors.New("condition ends with an opening parentheses")
		}
		var (
			rPos   = -1
			offset = lPos + 1
		)
		for {
			if offset >= len(tokens)-1 {
				return nil, errors.New("missing matching closing parentheses")
			}
			rPos = findToken(tokens[offset:], tokenTypeRPAR)
			if rPos < 0 {
				return nil, errors.New("missing matching closing parentheses")
			}
			rPos += offset
			if lPos+1 == rPos {
				return nil, errors.New("empty subexpression found")
			}
			if pos := findToken(tokens[offset:rPos], tokenTypeLPAR); pos < 0 {
				break
			}
			offset = rPos + 1
		}
		subNode, err := parse(tokens[lPos+1 : rPos])
		if err != nil {
			return nil, fmt.Errorf("could not parse subexpression: %s", err)
		}
		res = append(res[:lPos], append([]interface{}{subNode}, res[rPos+1:]...)...)
		tokens = append(tokens[:lPos], append([]token{{tokenType: tokenTypeNONE}}, tokens[rPos+1:]...)...)
	}
	for _, tokenType := range []tokenType{
		tokenTypeID,
		tokenTypeVAL,
		tokenTypeEQ,
		tokenTypeSET,
		tokenTypeNE,
		tokenTypeGT,
		tokenTypeGTE,
		tokenTypeLT,
		tokenTypeLTE,
		tokenTypeNOT,
		tokenTypeAND,
		tokenTypeOR,
	} {
		var tokenFound = true
		for tokenFound {
			tokenFound = false
			for i, elem := range res {
				token, _ := elem.(token)
				if token.tokenType != tokenType {
					continue
				}
				tokenFound = true
				switch tokenType {
				case tokenTypeID:
					res[i] = ID{valueNode{nodeValue: token.matched}}
				case tokenTypeVAL:
					res[i] = VAL{valueNode{nodeValue: token.matched}}
				default:
					if i+1 >= len(res) {
						return nil, fmt.Errorf("missing parameter for %s operator", tokenTypeString[tokenType])
					}
					subNode1, ok := res[i+1].(Node)
					if !ok {
						return nil, fmt.Errorf("parameter for %s operator is not a node (1), got: %s", tokenTypeString[tokenType], res[i+1])
					}
					res = append(res[:i+1], res[i+2:]...) // remove the (i+1)th element because it has become a sub node
					switch tokenType {
					case tokenTypeNOT:
						res[i] = NOT{oneSubNode{node: subNode1}}
					default:
						if i == 0 {
							return nil, fmt.Errorf("missing parameter for %s operator", tokenTypeString[tokenType])
						}
						subNode2, ok := res[i-1].(Node)
						if !ok {
							return nil, fmt.Errorf("parameter for %s operator is not a node (2)", tokenTypeString[tokenType])
						}
						n := twoSubNodes{subNode2, subNode1}
						switch tokenType {
						case tokenTypeEQ:
							res[i] = EQ{n}
						case tokenTypeSET:
							res[i] = SET{n}
						case tokenTypeNE:
							res[i] = NE{n}
						case tokenTypeGT:
							res[i] = GT{n}
						case tokenTypeGTE:
							res[i] = GTE{n}
						case tokenTypeLT:
							res[i] = LT{n}
						case tokenTypeLTE:
							res[i] = LTE{n}
						case tokenTypeAND:
							res[i] = AND{n}
						case tokenTypeOR:
							res[i] = OR{n}
						default:
							return nil, fmt.Errorf("invalid token type: %s", tokenTypeString[tokenType])
						}
						res = append(res[:i-1], res[i:]...) // remove the (i-1)the element because it has become a sub node
					}
					break
				}
			}
		}
	}
	if len(res) != 1 {
		return nil, errors.New("parse tree must have exactly one start node")
	}
	startNode, ok := res[0].(Node)
	if !ok {
		return nil, errors.New("start node is not a node")
	}
	return startNode, nil
}
