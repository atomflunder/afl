package runtime

import (
	"fmt"
	"strconv"
	"strings"

	"afl/src/parser"
)

type RuntimeValue struct {
	Type  string
	Value any
}

// EvaluateProgram evaluates a list of expressions (AST) in the given environment
func EvaluateProgram(ast []parser.Expression, env *Environment) RuntimeValue {
	lastEvaluated := RuntimeValue{Type: "null", Value: nil}

	for _, statement := range ast {
		lastEvaluated = Evaluate(statement, env)

		if lastEvaluated.Type == "return" {
			return lastEvaluated
		}
	}

	return lastEvaluated
}

// Evaluate is the main dispatcher that evaluates any expression
func Evaluate(expr parser.Expression, env *Environment) RuntimeValue {
	switch node := expr.(type) {
	case parser.NumberLiteralExpression:
		return evaluateNumberLiteral(node)
	case parser.StringLiteralExpression:
		return evaluateStringLiteral(node, env)
	case parser.IdentifierExpression:
		return evaluateIdentifier(node, env)
	case parser.BinaryExpression:
		return evaluateBinaryExpression(node, env)
	case parser.CallExpression:
		return evaluateCallExpression(node, env)
	case parser.VariableDeclaration:
		return evaluateVariableDeclaration(node, env)
	case parser.FunctionDeclaration:
		return evaluateFunctionDeclaration(node, env)
	case parser.IfStatement:
		return evaluateIfStatement(node, env)
	case parser.ForLoop:
		return evaluateForLoop(node, env)
	case parser.AssignmentExpression:
		return evaluateAssignmentExpression(node, env)
	case []parser.Property:
		return evaluateObjectLiteral(node, env)
	case parser.MemberExpression:
		return evaluateMemberExpression(node, env)
	case parser.InfinityExpression:
		return RuntimeValue{Type: "infinity", Value: nil}
	default:
		return RuntimeValue{Type: "null", Value: nil}
	}
}

func evaluateNumberLiteral(node parser.NumberLiteralExpression) RuntimeValue {
	value, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return RuntimeValue{Type: "null", Value: nil}
	}
	return RuntimeValue{Type: "number", Value: value}
}

func evaluateStringLiteral(node parser.StringLiteralExpression, env *Environment) RuntimeValue {
	interpolated := interpolateString(node.Value, env)
	return RuntimeValue{Type: "string", Value: interpolated}
}

// interpolateString replaces {variable} patterns with actual variable values
func interpolateString(str string, env *Environment) string {
	result := str
	i := 0

	for i < len(result) {
		// Find opening brace
		openIdx := strings.Index(result[i:], "{")
		if openIdx == -1 {
			break
		}

		openIdx += i

		// Find closing brace
		closeIdx := strings.Index(result[openIdx:], "}")
		if closeIdx == -1 {
			break
		}

		closeIdx += openIdx

		// Extract variable name
		varName := result[openIdx+1 : closeIdx]

		// Look up the variable
		val, err := env.getVariable(varName)
		if err == nil {
			// Convert value to string
			var strVal string
			if rv, ok := val.(RuntimeValue); ok {
				strVal = runtimeValueToString(rv)
			} else {
				strVal = fmt.Sprintf("%v", val)
			}

			// Replace the pattern with the value
			result = result[:openIdx] + strVal + result[closeIdx+1:]
			i = openIdx + len(strVal)
		} else {
			// Variable not found, move past this brace
			i = closeIdx + 1
		}
	}

	return result
}

func evaluateIdentifier(node parser.IdentifierExpression, env *Environment) RuntimeValue {
	val, err := env.getVariable(node.Symbol)
	if err != nil {
		// Undefined variable
		return RuntimeValue{Type: "null", Value: nil}
	}

	if rv, ok := val.(RuntimeValue); ok {
		return rv
	}

	return RuntimeValue{Type: "object", Value: val}
}

func evaluateBinaryExpression(node parser.BinaryExpression, env *Environment) RuntimeValue {
	if node.Left == nil {
		return evaluateUnaryExpression(node, env)
	}

	left := Evaluate(node.Left, env)
	right := Evaluate(node.Right, env)

	if left.Type == "number" && right.Type == "number" {
		return evaluateNumericBinary(left, right, node.Operator.Value)
	}

	if left.Type == "string" && right.Type == "string" {
		return evaluateStringBinary(left, right, node.Operator.Value)
	}

	if node.Operator.Value == "==" || node.Operator.Value == "!=" {
		return evaluateComparison(left, right, node.Operator.Value)
	}

	return RuntimeValue{Type: "null", Value: nil}
}

func evaluateNumericBinary(left, right RuntimeValue, operator string) RuntimeValue {
	leftVal := left.Value.(float64)
	rightVal := right.Value.(float64)

	var result float64
	switch operator {
	case "+":
		result = leftVal + rightVal
	case "-":
		result = leftVal - rightVal
	case "*":
		result = leftVal * rightVal
	case "/":
		if rightVal == 0 {
			return RuntimeValue{Type: "null", Value: nil}
		}
		result = leftVal / rightVal
	case "%":
		result = float64(int(leftVal) % int(rightVal))
	case "<":
		if leftVal < rightVal {
			return RuntimeValue{Type: "number", Value: 1.0}
		}
		return RuntimeValue{Type: "number", Value: 0.0}
	case ">":
		if leftVal > rightVal {
			return RuntimeValue{Type: "number", Value: 1.0}
		}
		return RuntimeValue{Type: "number", Value: 0.0}
	case "<=":
		if leftVal <= rightVal {
			return RuntimeValue{Type: "number", Value: 1.0}
		}
		return RuntimeValue{Type: "number", Value: 0.0}
	case ">=":
		if leftVal >= rightVal {
			return RuntimeValue{Type: "number", Value: 1.0}
		}
		return RuntimeValue{Type: "number", Value: 0.0}
	case "==":
		if leftVal == rightVal {
			return RuntimeValue{Type: "number", Value: 1.0}
		}
		return RuntimeValue{Type: "number", Value: 0.0}
	case "!=":
		if leftVal != rightVal {
			return RuntimeValue{Type: "number", Value: 1.0}
		}
		return RuntimeValue{Type: "number", Value: 0.0}
	default:
		return RuntimeValue{Type: "null", Value: nil}
	}
	return RuntimeValue{Type: "number", Value: result}
}

func evaluateStringBinary(left, right RuntimeValue, operator string) RuntimeValue {
	leftVal := left.Value.(string)
	rightVal := right.Value.(string)

	switch operator {
	case "+":
		return RuntimeValue{Type: "string", Value: leftVal + rightVal}
	default:
		return RuntimeValue{Type: "null", Value: nil}
	}
}

func evaluateComparison(left, right RuntimeValue, operator string) RuntimeValue {
	switch operator {
	case "==":
		if isEqual(left, right) {
			return RuntimeValue{Type: "number", Value: 1.0}
		}
		return RuntimeValue{Type: "number", Value: 0.0}
	case "!=":
		if !isEqual(left, right) {
			return RuntimeValue{Type: "number", Value: 1.0}
		}
		return RuntimeValue{Type: "number", Value: 0.0}
	default:
		return RuntimeValue{Type: "null", Value: nil}
	}
}

func isEqual(left, right RuntimeValue) bool {
	if left.Type != right.Type {
		return false
	}

	switch left.Type {
	case "number":
		return left.Value.(float64) == right.Value.(float64)
	case "string":
		return left.Value.(string) == right.Value.(string)
	case "null":
		return true
	default:
		return left.Value == right.Value
	}
}

func evaluateUnaryExpression(node parser.BinaryExpression, env *Environment) RuntimeValue {
	operand := Evaluate(node.Right, env)

	switch node.Operator.Value {
	case "!":
		if isTruthy(operand) {
			return RuntimeValue{Type: "number", Value: 0.0}
		}
		return RuntimeValue{Type: "number", Value: 1.0}
	case "-":
		// Numeric negation: -5
		if operand.Type == "number" {
			return RuntimeValue{Type: "number", Value: -(operand.Value.(float64))}
		}
		return RuntimeValue{Type: "null", Value: nil}
	default:
		return RuntimeValue{Type: "null", Value: nil}
	}
}

func isTruthy(val RuntimeValue) bool {
	switch val.Type {
	case "null":
		return false
	case "number":
		return val.Value.(float64) != 0
	case "string":
		return val.Value.(string) != ""
	default:
		return true
	}
}

func evaluateCallExpression(node parser.CallExpression, env *Environment) RuntimeValue {
	caller := Evaluate(node.Caller, env)

	args := make([]RuntimeValue, len(node.Args))
	for i, arg := range node.Args {
		args[i] = Evaluate(arg, env)
	}

	// Right now, print() is the only built in function
	if callerIdent, ok := node.Caller.(parser.IdentifierExpression); ok {
		switch callerIdent.Symbol {
		case "print":
			return callPrint(args)
		}
	}

	if caller.Type == "function" {
		fn := caller.Value.(*FunctionValue)
		return callUserFunction(fn, args)
	}

	return RuntimeValue{Type: "null", Value: nil}
}

func callPrint(args []RuntimeValue) RuntimeValue {
	var parts []string
	for _, arg := range args {
		parts = append(parts, runtimeValueToString(arg))
	}
	fmt.Println(strings.Join(parts, " "))
	return RuntimeValue{Type: "null", Value: nil}
}

func runtimeValueToString(val RuntimeValue) string {
	switch val.Type {
	case "number":
		num := val.Value.(float64)
		// Format as integer if it's a whole number
		if num == float64(int64(num)) {
			return fmt.Sprintf("%d", int64(num))
		}
		return fmt.Sprintf("%g", num)
	case "string":
		return val.Value.(string)
	case "null":
		return "null"
	default:
		return fmt.Sprintf("%v", val.Value)
	}
}

type FunctionValue struct {
	Parameters []string
	Body       []parser.Expression
	Env        *Environment
}

func evaluateFunctionDeclaration(node parser.FunctionDeclaration, env *Environment) RuntimeValue {
	fn := &FunctionValue{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}

	// Store the function in the environment
	env.declareVariable(node.Name, RuntimeValue{Type: "function", Value: fn}, false)

	return RuntimeValue{Type: "null", Value: nil}
}

func callUserFunction(fn *FunctionValue, args []RuntimeValue) RuntimeValue {
	scope := NewEnvironmentWithOuter(fn.Env)

	// Bind parameters to arguments
	for i, param := range fn.Parameters {
		if i < len(args) {
			scope.declareVariable(param, args[i], false)
		} else {
			scope.declareVariable(param, RuntimeValue{Type: "null", Value: nil}, false)
		}
	}

	result := RuntimeValue{Type: "null", Value: nil}
	for _, stmt := range fn.Body {
		result = Evaluate(stmt, scope)

		if result.Type == "return" {
			val := result.Value.(RuntimeValue)
			return val
		}
	}

	return result
}

func evaluateVariableDeclaration(node parser.VariableDeclaration, env *Environment) RuntimeValue {
	value := Evaluate(node.Value, env)
	isConstant := !node.IsMutable

	err := env.declareVariable(node.Name, value, isConstant)
	if err != nil {
		env.assignVariable(node.Name, value)
	}

	return value
}

func evaluateAssignmentExpression(node parser.AssignmentExpression, env *Environment) RuntimeValue {
	value := Evaluate(node.Right, env)

	if ident, ok := node.Left.(parser.IdentifierExpression); ok {
		env.assignVariable(ident.Symbol, value)
		return value
	}
	if member, ok := node.Left.(parser.MemberExpression); ok {
		obj := Evaluate(member.Object, env)
		if obj.Type == "object" {
			objMap := obj.Value.(map[string]RuntimeValue)
			objMap[member.Property] = value
			return value
		}
	}

	return RuntimeValue{Type: "null", Value: nil}
}

func evaluateIfStatement(node parser.IfStatement, env *Environment) RuntimeValue {
	condition := Evaluate(node.IfBranch.Condition, env)

	if isTruthy(condition) {
		return evaluateBlock(node.IfBranch.Body, env)
	}

	for _, branch := range node.ElseIfBranches {
		branchCondition := Evaluate(branch.Condition, env)
		if isTruthy(branchCondition) {
			return evaluateBlock(branch.Body, env)
		}
	}

	if len(node.ElseBranch.Body) > 0 {
		return evaluateBlock(node.ElseBranch.Body, env)
	}

	return RuntimeValue{Type: "null", Value: nil}
}

func evaluateForLoop(node parser.ForLoop, env *Environment) RuntimeValue {
	scope := NewEnvironmentWithOuter(env)

	start := Evaluate(node.Start, env)
	end := Evaluate(node.End, env)

	if end.Type == "infinity" {
		startVal := uint(0)
		if start.Type == "number" {
			startVal = uint(start.Value.(float64))
		}

		result := RuntimeValue{Type: "null", Value: nil}
		// Maximum value for a uint type.
		const maxIterations = ^uint(0)
		for i := startVal; i < startVal+maxIterations; i++ {
			scope.assignVariable(node.Iterator, RuntimeValue{Type: "number", Value: float64(i)})

			for _, stmt := range node.Body {
				result = Evaluate(stmt, scope)
				if result.Type == "return" {
					return result
				}
			}
		}
		return result
	}

	if start.Type != "number" || end.Type != "number" {
		return RuntimeValue{Type: "null", Value: nil}
	}

	startVal := int64(start.Value.(float64))
	endVal := int64(end.Value.(float64))

	result := RuntimeValue{Type: "null", Value: nil}
	for i := startVal; i < endVal; i++ {
		scope.assignVariable(node.Iterator, RuntimeValue{Type: "number", Value: float64(i)})

		for _, stmt := range node.Body {
			result = Evaluate(stmt, scope)
			if result.Type == "return" {
				return result
			}
		}
	}

	return result
}

func evaluateBlock(body []parser.Expression, env *Environment) RuntimeValue {
	scope := NewEnvironmentWithOuter(env)

	result := RuntimeValue{Type: "null", Value: nil}
	for _, stmt := range body {
		result = Evaluate(stmt, scope)
		if result.Type == "return" {
			return result
		}
	}

	return result
}

func evaluateObjectLiteral(properties []parser.Property, env *Environment) RuntimeValue {
	obj := make(map[string]RuntimeValue)
	for _, prop := range properties {
		obj[prop.Key] = Evaluate(prop.Value, env)
	}
	return RuntimeValue{Type: "object", Value: obj}
}

func evaluateMemberExpression(node parser.MemberExpression, env *Environment) RuntimeValue {
	obj := Evaluate(node.Object, env)

	if obj.Type == "object" {
		objMap := obj.Value.(map[string]RuntimeValue)
		if val, exists := objMap[node.Property]; exists {
			return val
		}
	}

	return RuntimeValue{Type: "null", Value: nil}
}
