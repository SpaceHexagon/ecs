// evaluator/evaluator.go
package evaluator

import (
	"fmt"
	"hash/fnv"

	"github.com/SpaceHexagon/ecs/util"

	"github.com/SpaceHexagon/ecs/ast"
	"github.com/SpaceHexagon/ecs/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func Eval(node ast.Node, env *object.Environment, objectContext *object.Hash) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env, objectContext)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env, objectContext)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.ClassStatement:
		val := Eval(node.Value, env, objectContext)
		if isError(val) {
			return val
		}

		h := fnv.New64a()
		h.Write([]byte(node.Name.Value))
		pair := val.(*object.Hash).Pairs[object.HashKey{Type: object.STRING_OBJ, Value: h.Sum64()}]
		constructor := pair.Value

		if constructor != nil {
			val.(*object.Hash).Constructor = constructor.(*object.Function)
		}
		env.Set(node.Name.Value, val)
		break
	case *ast.LetStatement:
		val := Eval(node.Value, env, objectContext)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.AssignmentStatement:
		val := Eval(node.Value, env, objectContext)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		break
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env, objectContext)
		// Expressions
	case *ast.PrefixExpression:
		right := Eval(node.Right, env, objectContext)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env, objectContext)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env, objectContext)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.CallExpression:
		function := Eval(node.Function, env, objectContext)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env, objectContext)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, objectContext)
	case *ast.IndexExpression:
		left := Eval(node.Left, env, objectContext)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env, objectContext)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env, objectContext)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env, objectContext)
	case *ast.IfExpression:
		return evalIfExpression(node, env, objectContext)
	case *ast.ForExpression:
		return evalForExpression(node, env, objectContext)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env, objectContext)
	case *ast.SleepExpression:
		return evalSleepExpression(node, env, objectContext)
	case *ast.NewExpression:
		return evalNewExpression(node, env, objectContext)
	// case *ast.ExecExpression:
	// 	return evalExecExpression(node, env, objectContext)
	case *ast.Identifier:
		return evalIdentifier(node, env, objectContext)

	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env, &object.Hash{})
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}
func evalStatements(stmts []ast.Statement, env *object.Environment, objectContext *object.Hash) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env, objectContext)
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment, objectContext *object.Hash) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env, objectContext)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}
func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
	objectContext *object.Hash,
) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env, objectContext)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}
func evalInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())

	}
}
func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}
func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalNewExpression(ne *ast.NewExpression, env *object.Environment, objectContext *object.Hash) object.Object {
	classData := evalIdentifier(ne.Name, env, objectContext)

	if classData.Type() != object.HASH_OBJ {
		return newError("new operator can only be used with Class or Hashmap. Invalid type: %s", classData.Type())
	}
	instance := util.CopyHashMap(classData)
	// if instance.Constructor == nil && instance.Pairs[ne.Name.Value] != nil {
	// 	instance.(*object.Hash).Constructor = instance.Pairs[ne.Name.Value].Value
	// 	instance.className = ne.Name.Value
	// 	bindContextToMethods(instance)
	// 	// if (instance.Pairs.builtin && instance.Pairs.builtin.Value) {
	// 	// 	extendBuiltinMethodEnvs(instance);
	// 	// }
	// }

	return instance
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment, objectContext *object.Hash) object.Object {
	condition := Eval(ie.Condition, env, objectContext)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env, objectContext)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env, objectContext)
	} else {
		return NULL
	}
}

func evalForExpression(fl *ast.ForExpression, env *object.Environment, objectContext *object.Hash) object.Object {
	var (
		index  int64 = 0
		length int64 = 0
	)
	rangeObj := Eval(fl.Range, env, objectContext)
	element := fl.Element.Value
	indexObj := &object.Integer{Value: index}

	if isError(rangeObj) {
		return rangeObj
	}
	rangeType := rangeObj.Type()
	var (
		err    object.Object
		result object.Object
	)

	if rangeType == object.INTEGER_OBJ {
		length := rangeObj.(*object.Integer).Value
		for index < length {
			indexObj.Value = index
			env.Set(element, indexObj)
			result = Eval(fl.Consequence, env, objectContext)
			if isError(result) {
				err = result
			}
			index++
		}
	} else if rangeType == object.ARRAY_OBJ {
		length = int64(len(rangeObj.(*object.Array).Elements))
		for index < length {
			indexObj.Value = index
			env.Set(element, indexObj)
			result = Eval(fl.Consequence, env, objectContext)
			if isError(result) {
				err = result
			}
			index++
		}
	} else if rangeType == object.STRING_OBJ {
		length = int64(len(rangeObj.(*object.String).Value))
		for index < length {
			indexObj.Value = index
			env.Set(element, indexObj)
			result = Eval(fl.Consequence, env, objectContext)
			if isError(result) {
				err = result
			}
			index++
		}
	} else if rangeType == object.HASH_OBJ {
		for _, v := range rangeObj.(*object.Hash).Pairs {
			env.Set(element, &object.String{Value: v.Value.Inspect()})
			result = Eval(fl.Consequence, env, objectContext)
			if isError(result) {
				err = result
			}
			index++
		}
		if err != nil {
			return newError("error in for loop %s", err)
		}
	} else {
		return newError("unknown range type in for loop: %s", rangeObj.Type())
	}

	return NULL
}

func evalWhileExpression(ie *ast.WhileExpression, env *object.Environment, objectContext *object.Hash) object.Object {
	condition := Eval(ie.Condition, env, objectContext)

	if isError(condition) {
		return condition
	}
	for isTruthy(condition) {
		Eval(ie.Consequence, env, objectContext)
		condition = Eval(ie.Condition, env, objectContext)
	}

	return NULL
}

func evalSleepExpression(se *ast.SleepExpression, env *object.Environment, objectContext *object.Hash) object.Object {
	duration := Eval(se.Duration, env, objectContext)
	if isError(duration) {
		return duration
	}
	// setTimeout(()=>{
	// 	Eval(se.Consequence, env, objectContext)
	// }, <number>duration.Inspect())

	return NULL
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return NULL
	}
	return arrayObject.Elements[idx]
}
func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}
func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
	objectContext *object.Hash,
) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}
func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
	objectContext *object.Hash,
) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env, objectContext)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env, objectContext)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func applyFunction(fn object.Object, args []object.Object, objectContext *object.Hash) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv, objectContext)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		// implement api context
		return fn.Fn(nil, nil, args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}
func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
