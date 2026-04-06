package parser

import "fmt"

type Parser struct {
	tokens []Token
}

// Returns a new parser instance of the given input string.
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

// Checks if it is not the end of the file (EOF) token.
func (p *Parser) notEof() bool {
	return p.tokens[0].Type != EOF
}

// Returns the token at the given index, without consuming it.
func (p *Parser) at(index int) Token {
	return p.tokens[index]
}

// Consumes and checks the next token to see if it matches the expected type.
func (p *Parser) expectType(tt TokenType) error {
	tokens := p.shiftTokens(1)
	token := tokens[0]

	if token.Type != tt {
		return fmt.Errorf("expected token of type %s, got %s", tt, token.Type)
	}
	return nil
}

// Shifts the tokens by the given number and returns the shifted tokens.
func (p *Parser) shiftTokens(by int) []Token {
	t := p.tokens[:by]
	p.tokens = p.tokens[by:]
	return t
}

// Skips comment tokens
func (p *Parser) skipComments() {
	for p.notEof() && (p.at(0).Type == LineComment || p.at(0).Type == BlockCommentStart || p.at(0).Type == BlockCommentEnd) {
		p.shiftTokens(1)
	}
}

// Parses a generic statement.
// Could be a variable declaration, function declaration, if statement, for loop, or an expression.
func (p *Parser) parseStatement() (Expression, error) {
	p.skipComments()

	switch p.at(0).Type {
	case Identifier:
		// Try to parse as variable declaration first
		result, err := p.parseVariableDeclaration()
		if err == nil {
			return result, nil
		}
		// If it's not a variable declaration, try as an expression
		return p.parseExpression()
	case Function:
		return p.parseFunctionDeclaration()
	case If:
		return p.parseIfStatement()
	case For:
		return p.parseForLoop()
	default:
		return p.parseExpression()
	}
}

type Expression interface{}

// Parses a general expression.
func (p *Parser) parseExpression() (Expression, error) {
	return p.parseAssignmentExpression()
}

type AssignmentExpression struct {
	Left  Expression
	Right Expression
}

// Parses an assignment expression.
func (p *Parser) parseAssignmentExpression() (Expression, error) {
	left, err := p.parseObjectExpression()
	if err != nil {
		return nil, err
	}

	if p.at(0).Type == Equals {
		p.shiftTokens(1)
		right, err := p.parseAssignmentExpression()
		if err != nil {
			return nil, err
		}

		return AssignmentExpression{
			Left:  left,
			Right: right,
		}, nil
	}

	return left, nil
}

type Property struct {
	Value Expression
	Key   string
}

// Parses an assignment expression for object literals, e.g. { x: 5, y: 10, }
func (p *Parser) parseObjectExpression() (Expression, error) {
	if p.at(0).Type != OpenCurly {
		return p.parseLogicalExpression()
	}

	p.shiftTokens(1)

	properties := []Property{}

	for p.notEof() && p.at(0).Type != CloseCurly {
		key := p.at(0)
		p.expectType(Identifier)

		p.expectType(Colon)

		value, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		properties = append(properties, Property{
			Key:   key.Value,
			Value: value,
		})

		// We want trailing commas to be required!
		p.expectType(Comma)
	}

	p.expectType(CloseCurly)

	return properties, nil
}

type BinaryExpression struct {
	Left     Expression
	Operator Token
	Right    Expression
}

// Parses any logical, comparison, additive, or multiplicative expression, based on operator precedence.
func (p *Parser) parseLogicalExpression() (Expression, error) {
	left, err := p.parseComparisonExpression()

	if err != nil {
		return BinaryExpression{}, err
	}

	for p.at(0).Type == LogicalAnd || p.at(0).Type == LogicalOr {
		operator := p.at(0)
		p.shiftTokens(1)

		right, err := p.parseComparisonExpression()
		if err != nil {
			return BinaryExpression{}, err
		}

		left = BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

// Handles comparison operators: ==, !=, <, >, <=, >=
func (p *Parser) parseComparisonExpression() (Expression, error) {
	left, err := p.parseAdditiveExpression()

	if err != nil {
		return nil, err
	}

	for p.at(0).Type == EqualsEquals || p.at(0).Type == NotEquals || p.at(0).Type == LessThan || p.at(0).Type == GreaterThan {
		operator := p.at(0)
		p.shiftTokens(1)

		right, err := p.parseAdditiveExpression()
		if err != nil {
			return BinaryExpression{}, err
		}

		left = BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

// Handles additive operators: +, -
func (p *Parser) parseAdditiveExpression() (Expression, error) {
	left, err := p.parseMultiplicativeExpression()
	if err != nil {
		return BinaryExpression{}, err
	}

	for p.at(0).Type == Plus || p.at(0).Type == Minus {
		operator := p.at(0)
		p.shiftTokens(1)

		right, err := p.parseMultiplicativeExpression()
		if err != nil {
			return BinaryExpression{}, err
		}

		left = BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

// Handles multiplicative operators: *, /, %
func (p *Parser) parseMultiplicativeExpression() (Expression, error) {
	left, err := p.parseUnaryExpression()
	if err != nil {
		return BinaryExpression{}, err
	}

	for p.at(0).Type == Asterisk || p.at(0).Type == Slash || p.at(0).Type == Percent {
		operator := p.at(0)
		p.shiftTokens(1)

		right, err := p.parseUnaryExpression()
		if err != nil {
			return BinaryExpression{}, err
		}

		left = BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left, nil
}

// Handles unary operators: !, -
func (p *Parser) parseUnaryExpression() (Expression, error) {
	if p.at(0).Type == ExclamationMark || p.at(0).Type == Minus {
		operator := p.at(0)
		p.shiftTokens(1)

		operand, err := p.parseUnaryExpression() // Recursive call to handle multiple unary operators in a row, e.g. "!!true" or "--5"
		if err != nil {
			return nil, err
		}

		return BinaryExpression{
			Left:     nil,
			Operator: operator,
			Right:    operand,
		}, nil
	}

	return p.parseCallMemberExpression()
}

// Parses member expressions and call expressions, e.g. foo.bar or foo()
func (p *Parser) parseCallMemberExpression() (Expression, error) {
	member, err := p.parseMemberExpression()
	if err != nil {
		return nil, err
	}

	if p.at(0).Type == OpenParen {
		return p.parseCallExpression(member)
	}

	return member, nil
}

type MemberExpression struct {
	Object   Expression
	Property string
}

// Parses member expressions, e.g. foo.bar.baz
func (p *Parser) parseMemberExpression() (Expression, error) {
	obj, err := p.parsePrimaryExpression()
	if err != nil {
		return MemberExpression{}, err
	}

	for p.at(0).Type == Period {
		p.shiftTokens(1)

		propertyToken := p.at(0)
		p.expectType(Identifier)

		obj = MemberExpression{
			Object:   obj,
			Property: propertyToken.Value,
		}

	}

	return obj, nil
}

type IdentifierExpression struct {
	Symbol string
}

type NumberLiteralExpression struct {
	Value string
}

type StringLiteralExpression struct {
	Value string
}

type InfinityExpression struct{}

// Parses primary expressions: identifiers, number literals, string literals, parenthesized expressions, and object literals.
func (p *Parser) parsePrimaryExpression() (Expression, error) {
	switch p.at(0).Type {
	case Identifier:
		identifier := p.at(0)
		p.shiftTokens(1)
		return IdentifierExpression{Symbol: identifier.Value}, nil
	case NumberLiteral:
		number := p.at(0)
		p.shiftTokens(1)
		return NumberLiteralExpression{Value: number.Value}, nil
	case StringLiteral:
		str := p.at(0)
		p.shiftTokens(1)
		return StringLiteralExpression{Value: str.Value}, nil
	case Infinity:
		p.shiftTokens(1)
		return InfinityExpression{}, nil
	default:
		return nil, fmt.Errorf("unexpected token type: %v %v", p.at(0).Type, p.tokens)
	}
}

// Parses a list of arguments in a function call, e.g. (x, y, z)
func (p *Parser) parseArgs() (Expression, error) {
	p.expectType(OpenParen)

	if p.at(0).Type == CloseParen {
		p.shiftTokens(1)
		return []Expression{}, nil
	}

	args, err := p.parseArgsList()
	if err != nil {
		return nil, err
	}

	p.expectType(CloseParen)

	return args, nil
}

// Parses the list of arguments in a function call, e.g. x, y, z
func (p *Parser) parseArgsList() ([]Expression, error) {
	args := []Expression{}

	for p.notEof() && p.at(0).Type == Comma {
		p.shiftTokens(0)
		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return args, nil
}

type CallExpression struct {
	Caller Expression
	Args   []Expression
}

// Parses a function call expression, e.g. foo(x, y, z)
func (p *Parser) parseCallExpression(caller Expression) (CallExpression, error) {
	err := p.expectType(OpenParen)
	if err != nil {
		return CallExpression{}, err
	}

	args := []Expression{}

	// Parse arguments if they exist
	for p.notEof() && p.at(0).Type != CloseParen {
		arg, err := p.parseExpression()
		if err != nil {
			return CallExpression{}, err
		}
		args = append(args, arg)

		// If there's a comma, there are more arguments
		if p.at(0).Type == Comma {
			p.shiftTokens(1)
		}
	}

	err = p.expectType(CloseParen)
	if err != nil {
		return CallExpression{}, err
	}

	callExpr := CallExpression{
		Caller: caller,
		Args:   args,
	}

	return callExpr, nil
}

type VariableDeclaration struct {
	IsMutable bool
	Name      string
	Value     Expression
}

// Parses variable declarations, e.g. x = 5, y: int = 6, a? = 10
func (p *Parser) parseVariableDeclaration() (VariableDeclaration, error) {
	if !p.notEof() || p.at(0).Type != Identifier {
		return VariableDeclaration{}, fmt.Errorf("not a variable declaration")
	}

	if len(p.tokens) < 2 {
		return VariableDeclaration{}, fmt.Errorf("not a variable declaration")
	}

	nextToken := p.at(1)

	if nextToken.Type != Equals && nextToken.Type != Colon && nextToken.Type != QuestionMark {
		return VariableDeclaration{}, fmt.Errorf("not a variable declaration")
	}

	nameToken := p.at(0)
	p.shiftTokens(1)
	name := nameToken.Value

	isMutable := false

	// Handle optional `?` for mutability and `:` for type annotation in any order
	// Valid patterns: `? :`, `: ?`, `?`, `:`
	if p.at(0).Type == QuestionMark {
		isMutable = true
		p.shiftTokens(1)
	}

	// Check for type annotation after mutability
	if p.at(0).Type == Colon {
		p.shiftTokens(1) // Skip the colon
		// Skip the type for now (int, float, string, bool)
		if p.notEof() {
			p.shiftTokens(1)
		}
	}

	err := p.expectType(Equals)
	if err != nil {
		return VariableDeclaration{}, err
	}

	value, err := p.parseExpression()
	if err != nil {
		return VariableDeclaration{}, err
	}

	if p.notEof() && p.at(0).Type == Semicolon {
		p.shiftTokens(1)
	}

	return VariableDeclaration{
		IsMutable: isMutable,
		Name:      name,
		Value:     value,
	}, nil
}

type FunctionDeclaration struct {
	Name       string
	Parameters []string
	Body       []Expression
}

// Parses function declarations, e.g. fun foo(x, y) { return x + y }
func (p *Parser) parseFunctionDeclaration() (FunctionDeclaration, error) {
	p.shiftTokens(1)

	nameToken := p.at(0)
	err := p.expectType(Identifier)
	if err != nil {
		return FunctionDeclaration{}, err
	}
	name := nameToken.Value

	params := p.parseParams()

	err = p.expectType(CloseParen)
	if err != nil {
		return FunctionDeclaration{}, err
	}

	body, err := p.parseFunctionBody()

	if err != nil {
		return FunctionDeclaration{}, err
	}

	fn := FunctionDeclaration{
		Name:       name,
		Parameters: params,
		Body:       body,
	}

	return fn, nil
}

func (p *Parser) parseParams() []string {
	err := p.expectType(OpenParen)
	if err != nil {
		return []string{}
	}

	if p.notEof() && p.at(0).Type == CloseParen {
		// The caller consumes the closing parenthesis, so we don't shift it here
		return []string{}
	}

	params := []string{}

	expr, err := p.parseExpression()
	if err != nil {
		return params
	}

	identExpr, ok := expr.(IdentifierExpression)
	if !ok {
		return params
	}

	params = append(params, identExpr.Symbol)

	for p.notEof() && p.at(0).Type == Comma {
		p.shiftTokens(1)

		expr, err := p.parseExpression()
		if err != nil {
			return params
		}

		identExpr, ok := expr.(IdentifierExpression)
		if !ok {
			return params
		}

		params = append(params, identExpr.Symbol)
	}

	return params
}

// Parses the body of a function declaration, which is a block of expressions enclosed in curly braces.
func (p *Parser) parseFunctionBody() ([]Expression, error) {
	err := p.expectType(OpenCurly)
	if err != nil {
		return nil, err
	}

	body := []Expression{}

	for p.notEof() && p.at(0).Type != CloseCurly {
		if p.at(0).Type == Return {
			p.shiftTokens(1) // Skip return keyword
			returnValue, err := p.parseExpression()
			if err != nil {
				return nil, err
			}
			// Skip optional semicolon
			if p.at(0).Type == Semicolon {
				p.shiftTokens(1)
			}
			body = append(body, returnValue)
			continue
		}

		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		// Skip optional semicolon
		if p.at(0).Type == Semicolon {
			p.shiftTokens(1)
		}

		body = append(body, expr)
	}

	err = p.expectType(CloseCurly)
	if err != nil {
		return nil, err
	}

	return body, nil
}

type Branch struct {
	Condition Expression
	Body      []Expression
}

type IfStatement struct {
	Condition      Expression
	IfBranch       Branch
	ElseIfBranches []Branch
	ElseBranch     Branch
}

// Parses if statements, including else if and else branches.
func (p *Parser) parseIfStatement() (IfStatement, error) {

	p.shiftTokens(1)
	err := p.expectType(OpenParen)
	if err != nil {
		return IfStatement{}, err
	}

	condition, err := p.parseExpression()
	if err != nil {
		return IfStatement{}, err
	}

	err = p.expectType(CloseParen)
	if err != nil {
		return IfStatement{}, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return IfStatement{}, err
	}

	elseIfBranches := []Branch{}

	for p.at(0).Type == ElseIf {
		p.shiftTokens(1)
		err := p.expectType(OpenParen)
		if err != nil {
			return IfStatement{}, err
		}

		// Same as the initial condition, but for the else if branch
		elseIfCondition, err := p.parseExpression()
		if err != nil {
			return IfStatement{}, err
		}

		err = p.expectType(CloseParen)
		if err != nil {
			return IfStatement{}, err
		}

		elseIfBody, err := p.parseBlock()
		if err != nil {
			return IfStatement{}, err
		}

		elseIfBranch := Branch{
			Condition: elseIfCondition,
			Body:      elseIfBody,
		}

		elseIfBranches = append(elseIfBranches, elseIfBranch)
	}

	var elseBranch Branch
	if p.at(0).Type == Else {
		p.shiftTokens(1)

		elseBody, err := p.parseBlock()
		if err != nil {
			return IfStatement{}, err
		}

		elseBranch = Branch{
			Body: elseBody,
		}
	}

	return IfStatement{
		Condition:      condition,
		IfBranch:       Branch{Condition: condition, Body: body},
		ElseIfBranches: elseIfBranches,
		ElseBranch:     elseBranch,
	}, nil
}

func (p *Parser) parseBlock() ([]Expression, error) {
	err := p.expectType(OpenCurly)
	if err != nil {
		return nil, err
	}

	expressions := []Expression{}
	for p.notEof() && p.at(0).Type != CloseCurly {
		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
	}
	err = p.expectType(CloseCurly)
	if err != nil {
		return nil, err
	}
	return expressions, nil

}

type ForLoop struct {
	Iterator string
	Start    Expression
	End      Expression
	Body     []Expression
}

// Parses for loops, e.g. for (i in 0 -> 10) { ... }
func (p *Parser) parseForLoop() (ForLoop, error) {
	p.shiftTokens(1)

	err := p.expectType(OpenParen)
	if err != nil {
		return ForLoop{}, err
	}

	iteratorToken := p.at(0)
	err = p.expectType(Identifier)
	if err != nil {
		return ForLoop{}, err
	}
	iterator := iteratorToken.Value

	err = p.expectType(In)
	if err != nil {
		return ForLoop{}, err
	}

	start, err := p.parseExpression()
	if err != nil {
		return ForLoop{}, err
	}

	err = p.expectType(ArrowRight)
	if err != nil {
		return ForLoop{}, err
	}

	end, err := p.parseExpression()
	if err != nil {
		return ForLoop{}, err
	}

	err = p.expectType(CloseParen)
	if err != nil {
		return ForLoop{}, err
	}

	body, err := p.parseBlock()
	if err != nil {
		return ForLoop{}, err
	}

	return ForLoop{
		Iterator: iterator,
		Start:    start,
		End:      end,
		Body:     body,
	}, nil
}

func (p *Parser) GetAst() ([]Expression, error) {
	program := []Expression{}

	for p.notEof() {
		// Skip optional semicolons between statements
		for p.notEof() && p.at(0).Type == Semicolon {
			p.shiftTokens(1)
		}

		if !p.notEof() {
			break
		}

		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}

		program = append(program, statement)
	}

	return program, nil
}
