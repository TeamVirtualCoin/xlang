package main
import(
	"fmt"
	"regexp"
)
//Tokens
type token struct {
	Type string
	Value string
}
var reserveds = []token{token{Type : "int",Value : "int"},
	token{Type : "string",Value : "str"},
	token{Type : "bool",Value : "bool"},
	token{Type : "block",Value : "block"},
	token{Type : "var",Value : "var"},
	token{Type : "if",Value : "if"},
}
var literals = []token{token{Type : "oparen",Value : "("},
	token{Type : "cparen",Value : ")"},
	token{Type : "obrace",Value : "{"},
	token{Type : "cbrace",Value : "}"},
	token{Type : "assign",Value : ":"},
	token{Type : "dic",Value : "\""},
	token{Type : "ic",Value : "'"},
	token{Type : "assign",Value : ":"},
	token{Type : "plus",Value : "+"},
	token{Type : "minus",Value : "-"},
	token{Type : "multiply",Value : "*"},
	token{Type : "divide",Value : "/"},
	token{Type : "break",Value : "?"},
	token{Type : "equal",Value : "="},
	token{Type : "lte",Value : "<="},
	token{Type : "gte",Value : ">="},
	token{Type : "lt",Value : "<"},
	token{Type : "gt",Value : ">"},
}
//Lexer
func lexer(input string) []token {
	var output []token
	var current = 0
	input += "\n"
	for current < len(input) {
		char := string(input[current])
		if isNumber(char) {
			current++
			value := char
			char = string(input[current])
			for isNumber(char) {
				value += char
				current++
				char = string(input[current])
			}
			output = append(output,token{Type : "number",Value : value})
			continue
		} else if isString(char) {
			current++
			value := char
			char = string(input[current])
			for isString(char) {
				value += char
				current++
				char = string(input[current])
			}
			isRes := false
			for q := range reserveds {
				thisres := reserveds[q]
				if value == thisres.Value {
					output = append(output,token{Type : thisres.Type,Value : thisres.Value})
					isRes = true
				}
			}
			if isRes == false {
				output = append(output,token{Type : "name",Value : value})
			}
			continue
		} else if isLiteral(char) {
			for k := range literals {
				this := literals[k]
				if len(this.Value) == 2 && char + string(input[current + 1]) == this.Value {
					output = append(output,token{Type : this.Type,Value : this.Value})
					current += 2
					break
				} else if char == this.Value {
					output = append(output,token{Type : this.Type,Value : this.Value})
					current++
					break
				}
				continue
			}
			continue
		} else if char == " " {
			current++
			continue
		} else if isLine(char) {
			current++
			continue
		} else {
			current++
			continue
		}
		break
	}
	return output
}
func isNumber(char string) bool {
	isNum, _ := regexp.MatchString("\\d",char)
	return isNum
}
func isString(char string) bool {
	isStr, _ := regexp.MatchString("\\w",char)
	return isStr
}
func isLiteral(char string) bool {
	if isNumber(char) == false && isString(char) == false && char != " " {
		if char != "\n" && char != "\t" {return true} else {return false}
	}
	return false
}
func isLine(char string) bool {
	if char == "\n" {return true}
	return false
}
//Parser
type node struct {
	NodeType string
	Type string
	Value string
	Body []node
	Param []node
}
var lc = 0
var pc = 0
var pt []token
func parser(tokens []token) node {
	pt = tokens
	pt = append(pt,token{Type : "break"})
	program := node{NodeType : "Program",Body : []node{}}
	for lc = range pt {
		if pc != lc {
			continue
		} else {
			parse := parse()
			if parse.NodeType != "ignore" {program.Body = append(program.Body,parse)}
		}
	}
	return program
}
func parse() node {
	token := pt[pc]
	pc++
	if token.Type == "number" {
		return node{NodeType : "NumberLiteral",Type : "number",Value : token.Value}
	} else if token.Type == "name" {
		return node{NodeType : "NameLiteral",Type : "name",Value : token.Value}
	} else if token.Type == "plus" {
		return node{NodeType : "Operate",Type : "plus",Param : []node{}}
	} else if token.Type == "minus" {
		return node{NodeType : "Operate",Type : "minus"}
	} else if token.Type == "mutiply" {
		return node{NodeType : "Operate",Type : "multiply"}
	} else if token.Type == "divide" {
		return node{NodeType : "Operate",Type : "Divide"}
	} else if token.Type == "lt" {
		return node{NodeType : "Comp",Type : "lt"}
	} else if token.Type == "gt" {
		return node{NodeType : "Comp",Type : "gt"}
	} else if token.Type == "lte" {
		return node{NodeType : "Comp",Type : "lte"}
	} else if token.Type == "gte" {
		return node{NodeType : "Comp",Type : "gte"}
	} else if token.Type == "equal" {
		return node{NodeType : "Comp",Type : "equal"}
	} else if token.Type == "obrace" {
		expr := node{NodeType : "BraceExp",Body : []node{}}
		for _,k := range pt[pc:] {
			if k.Type == "cbrace" {
				break
			} else {
				parse := parse()
				if parse.NodeType != "ignore" {expr.Body = append(expr.Body,parse)}
				continue
			}
			break
		}
		return expr
	} else if token.Type == "oparen" {
		expr := node{NodeType : "ParenExp",Body : []node{}}
		for _,k := range pt[pc:] {
			if k.Type == "cparen" {
				break
			} else {
				parse := parse()
				if parse.NodeType != "ignore" {expr.Body = append(expr.Body,parse)}
				continue
			}
			break
		}
		return expr
	} else if token.Type == "cbrace" || token.Type == "cparen" || token.Type == "break" || token.Type == "assign" || token.Type == "ignore" {
		return node{NodeType : "ignore"}
	} else if token.Type == "var" || token.Type == "int" || token.Type == "bool" || token.Type == "str"{
		var expr node
		if token.Type == "var" {expr = node{NodeType : "VariableExp"}
		} else if token.Type == "int" {expr = node{NodeType : "IntegerExp"}
		} else if token.Type == "bool" {expr = node{NodeType : "BooleanExp"}
		} else if token.Type == "string" {expr = node{NodeType : "StringExp"}}
		for _,k := range pt[pc:] {
			if k.Type == "break" {
				break
			} else {
				parse := parse()
				if parse.NodeType != "ignore" {expr.Param = append(expr.Param,parse)}
			}
		}
		return expr
	} else if token.Type == "block" {
		expr := node{NodeType : "block",Param : []node{},Body : []node{}}
		for _,k := range pt[pc:] {
			if k.Type == "cbrace" {
				break
			} else {
				parse := parse()
				if parse.NodeType != "ignore" {expr.Param = append(expr.Body,parse)}
			}
		}
	} else if token.Type == "dic" || token.Type == "ic" {
		expr := node{NodeType : "StringExp"}
		for _,k := range pt[pc:] {
			if k.Type == token.Type {
				pc++
				break
			} else {
				parse := parse()
				if parse.NodeType != "ignore"{expr.Body = append(expr.Body,parse)}
				continue
			}
			break
		}
		return expr
	}
	return node{}
}
func main() {
	a := lexer(`
	block a() {
		var a : "test"?
	}?`)
	fmt.Println(parser(a))
}
