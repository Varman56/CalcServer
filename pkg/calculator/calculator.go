package calculator

import (
	"errors"
	"strconv"
	"unicode"
)

var (
	ErrDivisionByZero           = errors.New("division by zero")
	ErrIncorrectBracketSequence = errors.New("incorrect count of brackets")
	ErrMultipleOperands         = errors.New("multiple operands in a row")
	ErrInvalidExpression        = errors.New("invalid expression")
	ErrConvertingToFloat64      = errors.New("failure to convert to float64")
	ErrUndefinedOperand         = errors.New("undefined operand")
)

type Token struct {
	is_num  bool
	operand rune
	num     float64
}

type Expression interface {
	ParsingExpression()
	CalculateExpression([]Token) (int, interface{})
}

type Arithmetic struct {
	expression            string
	parsed_expression     []Token
	result                float64
	is_invalid_expression bool
	err                   error
}

func (e *Arithmetic) setError(err error) {
	e.is_invalid_expression = true
	e.err = err
}

func isOperand(char rune) bool {
	switch char {
	case '+':
		return true
	case '-':
		return true
	case '*':
		return true
	case '/':
		return true
	default:
		return false
	}
}

func (e *Arithmetic) ParsingExpression() {
	var brackets_balance int = 0
	var result_expression []Token
	var index int = 0
	var last_digit_index int = 0
	expr := []rune(e.expression)
	for ; index < len(expr); index++ {
		var symbol rune = expr[index]
		switch {
		case symbol == '(':
			if len(result_expression) > 0 && !result_expression[len(result_expression)-1].is_num && result_expression[len(result_expression)-1].operand == ')' {
				result_expression = append(result_expression, Token{operand: '*'})
			}
			brackets_balance++
			result_expression = append(result_expression, Token{operand: '('})
		case symbol == ')':
			if brackets_balance--; brackets_balance < 0 {
				e.setError(ErrIncorrectBracketSequence)
				return
			}
			result_expression = append(result_expression, Token{operand: ')'})
		case isOperand(symbol):
			if len(result_expression) > 0 && !result_expression[len(result_expression)-1].is_num && isOperand(result_expression[len(result_expression)-1].operand) {
				e.setError(ErrMultipleOperands)
				return
			}
			result_expression = append(result_expression, Token{operand: symbol})
		case unicode.IsDigit(symbol):
			last_digit_index = index + 1
			for last_digit_index < len(expr) && (unicode.IsDigit(expr[last_digit_index]) || expr[last_digit_index] == '.') {
				last_digit_index++
			}
			num, err := strconv.ParseFloat(e.expression[index:last_digit_index], 64)
			if err != nil {
				e.setError(ErrConvertingToFloat64)
				return
			}
			if len(result_expression) > 0 && !result_expression[len(result_expression)-1].is_num && result_expression[len(result_expression)-1].operand == ')' {
				result_expression = append(result_expression, Token{operand: '*'})
			}
			result_expression = append(result_expression, Token{is_num: true, num: num})
			index = last_digit_index - 1
		default:
			e.setError(ErrUndefinedOperand)
			return
		}
	}
	if brackets_balance != 0 {
		e.setError(ErrIncorrectBracketSequence)
		return
	}
	e.parsed_expression = result_expression
}

func makeOperations(cleared_tokens []Token, high_priority bool) ([]Token, error) {
	var result_tokens []Token
	if high_priority {
		var index int = 0
		for ; index < len(cleared_tokens); index++ {
			token := cleared_tokens[index]
			switch token.operand {
			case '*':
				if len(result_tokens) == 0 || !result_tokens[len(result_tokens)-1].is_num ||
					index == len(cleared_tokens)-1 || !cleared_tokens[index+1].is_num {
					return cleared_tokens, ErrMultipleOperands
				}
				result_tokens[len(result_tokens)-1].num *= cleared_tokens[index+1].num
				index += 1
			case '/':
				if len(result_tokens) == 0 || !result_tokens[len(result_tokens)-1].is_num ||
					index == len(cleared_tokens)-1 || !cleared_tokens[index+1].is_num {
					return cleared_tokens, ErrMultipleOperands
				}
				if cleared_tokens[index+1].num == 0 {
					return cleared_tokens, ErrDivisionByZero
				}
				result_tokens[len(result_tokens)-1].num /= cleared_tokens[index+1].num
				index += 1
			default:
				result_tokens = append(result_tokens, cleared_tokens[index])
			}
		}
		return result_tokens, nil
	}
	var next_coef float64 = 1
	for _, token := range cleared_tokens {
		switch {
		case token.operand == '-':
			next_coef *= -1
		case token.is_num:
			result_tokens = append(result_tokens, Token{is_num: true, num: token.num * next_coef})
			next_coef = 1
		case token.operand == '+':
			continue
		default:
			return cleared_tokens, ErrUndefinedOperand
		}
	}
	result := make([]Token, 1)
	result[0].is_num = true
	for _, token := range result_tokens {
		result[0].num += token.num
	}
	return result, nil
}

func (e *Arithmetic) CalculateExpression(expr []Token) (int, interface{}) {
	if len(expr) == 0 {
		e.setError(ErrInvalidExpression)
		return 0, 0
	}
	var cleared_expr []Token
	var index int = 0
	var break_condition bool = false
	for ; index < len(expr); index++ {
		switch expr[index].operand {
		case '(':
			jump_index, inside_expr := e.CalculateExpression(expr[index+1:])
			inside_float, ok := inside_expr.(float64)
			if !ok {
				e.setError(ErrInvalidExpression)
				return 0, 0
			}
			cleared_expr = append(cleared_expr, Token{num: inside_float, is_num: true})
			index += jump_index
		case ')':
			break_condition = true
		default:
			cleared_expr = append(cleared_expr, expr[index])
		}
		if break_condition {
			break
		}
	}
	if e.is_invalid_expression {
		return 0, 0
	}
	high_priority_calced, err := makeOperations(cleared_expr, true)
	if err != nil {
		e.setError(err)
		return 0, 0
	}
	res, err := makeOperations(high_priority_calced, false)
	if err != nil {
		e.setError(err)
		return 0, 0
	}
	e.result = res[0].num
	return index + 1, res[0].num
}

func Calc(expression string) (float64, error) {
	arithmetic := Arithmetic{expression: expression}
	arithmetic.ParsingExpression()
	if arithmetic.is_invalid_expression {
		return 0, arithmetic.err
	}
	arithmetic.CalculateExpression(arithmetic.parsed_expression)
	if arithmetic.is_invalid_expression {
		return 0, arithmetic.err
	}
	return arithmetic.result, nil
}
