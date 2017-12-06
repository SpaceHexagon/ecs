package parser

import (
	"fmt"

	"github.com/SpaceHexagon/ecs/ast"
	"github.com/SpaceHexagon/ecs/lexer"
	"github.com/SpaceHexagon/ecs/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(expectedType token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", expectedType, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// mvp implementation to satisfy our tests
	// keep skipping until semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// func (p *Parser) parseIdentifier() {
// 	identifier = newIdentifierASTNode()
// 	identifier.token = currentToken()
// 	return identifier
// }

// func (p *Parser) parseExpression() {
// 	if currentToken() == INTEGER_TOKEN {
// 		if nextToken() == PLUS_TOKEN {
// 			return parseOperatorExpression()
// 		} else if nextToken() == SEMICOLON_TOKEN {
// 			return parseIntegerLiteral()
// 		}
// 	} else if currentToken() == LEFT_PAREN {
// 		return parseGroupedExpression()
// 	}
// 	// [...]
// }

// func (p *Parser) parseOperatorExpression() {
// 	operatorExpression = newOperatorExpression()
// 	operatorExpression.left = parseIntegerLiteral()
// 	operatorExpression.operator = currentToken()
// 	operatorExpression.right = parseExpression()
// 	return operatorExpression()
// }
