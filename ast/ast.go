package ast

import (
	"bytes"
	"strings"

	"github.com/SpaceHexagon/ecs/token"
)

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

type ForExpression struct {
	Token       token.Token // The 'for' token
	Element     Identifier
	Range       Expression
	Consequence *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for")
	out.WriteString(fe.Element.String())
	out.WriteString(fe.Range.String())
	out.WriteString(" ")
	out.WriteString(fe.Consequence.String())
	return out.String()
}

type WhileExpression struct {
	Token       token.Token // The 'while' token
	Condition   Expression
	Consequence *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Consequence.String())
	return out.String()
}

type SleepExpression struct {
	Token       token.Token // The 'sleep' token
	Duration    Expression
	Consequence *BlockStatement
}

func (se *SleepExpression) expressionNode()      {}
func (se *SleepExpression) TokenLiteral() string { return se.Token.Literal }
func (se *SleepExpression) String() string {
	var out bytes.Buffer
	out.WriteString("sleep")
	out.WriteString(se.Duration.String())
	out.WriteString(" ")
	out.WriteString(se.Consequence.String())
	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type NewExpression struct {
	Token token.Token // the token.NEW token
	Name  *Identifier
}

func (ce *NewExpression) expressionNode() {}
func (ce *NewExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *NewExpression) String() string {
	out := ce.TokenLiteral() + " " + ce.Name.String()

	return out
}

type ExecExpression struct {
	Token token.Token // the token.NEW token
	Name  Expression
}

func (ce *ExecExpression) expressionNode() {}
func (ce *ExecExpression) TokenLiteral() string {
	return ce.Token.Literal
}
func (ce *ExecExpression) String() string {
	out := ce.TokenLiteral() + " " + ce.Name.String()

	return out
}

type AssignmentStatement struct {
	Name  Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode() {}
func (as *AssignmentStatement) TokenLiteral() string {
	return as.Name.Value
}
func (as *AssignmentStatement) String() string {
	var out string = ""

	out += as.Name.String()
	out += " = "
	if as.Value != nil {
		out += as.Value.String()
	}
	out += ";"
	return out
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

/*
export class ClassStatement implements Statement {
	constructor(
		public Token:  token.Token, // the token.LET token
		public Name?:  Identifier,
		public Value?: HashLiteral,
		public NodeName: string = "ClassStatement"
	) {
	}
	public statementNode() {}
	public TokenLiteral(): string {
		return this.Token.Literal
	}
	public String(): string {
		var out:string
		out += this.TokenLiteral() + " "
		out += this.Name.String()
		out += " = "
		if (this.Value != null) {
			out += this.Value.String()
		}
		out += ";"
		return out;
	}
}
*/

type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type IndexAssignmentExpression struct {
	Token      token.Token // The [ token
	Left       Expression
	Index      Expression
	Assignment Expression
}

func (iae *IndexAssignmentExpression) expressionNode()      {}
func (iae *IndexAssignmentExpression) TokenLiteral() string { return iae.Token.Literal }
func (iae *IndexAssignmentExpression) String() string {
	out := ""
	out += iae.Left.String()
	out += "."
	out += iae.Index.String()
	out += "="
	out += iae.Assignment.String()
	return out
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type StringLiteral struct {
	Token token.Token
	Value string
	Node
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	return out.String()
}

func newProgramASTNode() *Program {
	p := Program{}
	return &p
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
