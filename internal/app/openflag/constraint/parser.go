package constraint

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/golang-collections/collections/stack"
)

var (
	// ErrFailedToParseExpression represents an error that we return when we can not parse the given expression.
	ErrFailedToParseExpression = errors.New("failed to parse expression")
)

// Parser represents a golang parser for parsing the main constraint based on the given expression.
type Parser struct{}

func (p Parser) isOperand(c string) bool {
	return "A" <= c && c <= "Z"
}

func (p Parser) isOperator(c string) bool {
	return c == NotConstraintName ||
		c == IntersectionConstraintName ||
		c == UnionConstraintName
}

func (p Parser) isParentheses(c string) bool {
	return c == "(" || c == ")"
}

func (p Parser) isValidCharacter(c string) bool {
	return p.isOperand(c) || p.isOperator(c) || p.isParentheses(c)
}

func (p Parser) findConstraint(identifier string, constraints map[string]model.Constraint) *model.Constraint {
	c, ok := constraints[identifier]
	if ok {
		return &c
	}

	return nil
}

func (p Parser) precedence(c string) int {
	table := map[string]int{
		"(":                        0,
		UnionConstraintName:        1,
		IntersectionConstraintName: 2,
		NotConstraintName:          3,
	}

	return table[c]
}

func (p Parser) generateOperator(operator string, args []model.Constraint) (*model.Constraint, error) {
	switch operator {
	case IntersectionConstraintName:
		p := IntersectionConstraint{
			Constraints: args,
		}

		parameters, err := json.Marshal(&p)
		if err != nil {
			return nil, err
		}

		return &model.Constraint{
			Name:       operator,
			Parameters: parameters,
		}, nil
	case UnionConstraintName:
		p := UnionConstraint{
			Constraints: args,
		}

		parameters, err := json.Marshal(&p)
		if err != nil {
			return nil, err
		}

		return &model.Constraint{
			Name:       operator,
			Parameters: parameters,
		}, nil
	case NotConstraintName:
		p := NotConstraint{
			Constraint: args[0],
		}

		parameters, err := json.Marshal(&p)
		if err != nil {
			return nil, err
		}

		return &model.Constraint{
			Name:       operator,
			Parameters: parameters,
		}, nil
	}

	return nil, errors.New("invalid expression operator")
}

// Parse parses and returns the main constraint based on the given expression.
// nolint:funlen
func (p Parser) Parse(expression string, constraints map[string]model.Constraint) (*model.Constraint, error) {
	exp := ""

	for i := 0; i < len([]rune(expression)); i++ {
		char := string([]rune(expression)[i])

		if char == " " || char == "\n" || char == "\r" || char == "\t" {
			continue
		}

		if !p.isValidCharacter(char) {
			return nil, fmt.Errorf("invalid character %s in expression", char)
		}

		if p.isOperand(char) && p.findConstraint(char, constraints) == nil {
			return nil, fmt.Errorf("undefined operand %s in expression", char)
		}

		exp += char
	}

	if exp == "" {
		return nil, errors.New("at least one operand should be in expression")
	}

	postfix := ""
	exp = "(" + exp + ")"

	parserStack := stack.New()
	i := 0

	for i < len([]rune(exp)) {
		char := string([]rune(exp)[i])

		if p.isOperand(char) {
			postfix += char
			i++
		} else if char == "(" {
			parserStack.Push(char)
			i++
		} else if p.isOperator(char) {
			if parserStack.Peek() != nil && p.precedence(char) > p.precedence(parserStack.Peek().(string)) {
				parserStack.Push(char)
				i++
			} else {
				postfix += parserStack.Peek().(string)
				parserStack.Pop()
			}
		} else if char == ")" {
			for parserStack.Peek() != nil && parserStack.Peek().(string) != "(" {
				postfix += parserStack.Peek().(string)
				parserStack.Pop()
			}

			parserStack.Pop()

			i++
		} else {
			return nil, ErrFailedToParseExpression
		}
	}

	parserStack = stack.New()

	for i := 0; i < len([]rune(postfix)); i++ {
		char := string([]rune(postfix)[i])

		if p.isOperand(char) {
			parserStack.Push(p.findConstraint(char, constraints))
		} else if char == IntersectionConstraintName || char == UnionConstraintName {
			args := []model.Constraint{}

			arg0 := parserStack.Pop()
			arg1 := parserStack.Pop()
			if arg0 == nil || arg1 == nil {
				return nil, ErrFailedToParseExpression
			}

			args = append(args, *arg0.(*model.Constraint))
			args = append(args, *arg1.(*model.Constraint))

			c, err := p.generateOperator(char, args)
			if err != nil {
				return nil, err
			}

			parserStack.Push(c)
		} else if char == NotConstraintName {
			args := []model.Constraint{}

			arg := parserStack.Pop()
			if arg == nil {
				return nil, ErrFailedToParseExpression
			}

			args = append(args, *arg.(*model.Constraint))

			c, err := p.generateOperator(char, args)
			if err != nil {
				return nil, err
			}

			parserStack.Push(c)
		} else {
			return nil, ErrFailedToParseExpression
		}
	}

	if parserStack.Len() == 1 {
		return parserStack.Pop().(*model.Constraint), nil
	}

	return nil, ErrFailedToParseExpression
}
