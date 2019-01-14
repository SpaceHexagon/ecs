package util

import (
	"github.com/SpaceHexagon/ecs/object"
)

// CopyObject returns a copy of any primitive object
func CopyObject(valueNode object.Object) object.Object {
	switch valueNode.Type() {
	case "BOOLEAN":
		return &object.Boolean{Value: valueNode.(*object.Boolean).Value}
	case "INTEGER":
		return &object.Integer{Value: valueNode.(*object.Integer).Value}
	case "FLOAT":
		return &object.Float{Value: valueNode.(*object.Float).Value}
	case "STRING":
		return &object.String{Value: valueNode.(*object.String).Value}
	case "ARRAY":
		return CopyArray(valueNode.(*object.Array))
	case "HASH":
		return CopyHashMap(valueNode)
	case object.FUNCTION_OBJ:
	case object.BUILTIN_OBJ:
	default:
		return &object.Null{}
	}
	return &object.Null{}
}

// CopyArray returns a deep copy of an existing object.Array
func CopyArray(array *object.Array) object.Object {
	var (
		elements []object.Object
	)
	for _, elem := range array.Elements {
		newObj := CopyObject(elem)
		elements = append(elements, newObj)
	}
	return &object.Array{Elements: elements}
}

// CopyHashMap creates a new object.Hash with the values of an existing one
// static fields and constructors are copied by reference
func CopyHashMap(data object.Object) object.Object {
	pairData := data.(*object.Hash).Pairs

	pairs := make(map[object.HashKey]object.HashPair)

	for key, pair := range pairData {
		valueNode := pair.Value
		keyNode := pair.Key
		isStatic := valueNode.Type() == "FUNCTION" || pair.Modifiers != nil && hasModifier(pair.Modifiers, 1)
		var (
			newPair object.HashPair
		)
		if isStatic {
			pairs[key] = pair
		} else {
			NewValue := CopyObject(valueNode)
			newPair = object.HashPair{Key: keyNode, Value: NewValue}
			if pair.Modifiers != nil {
				newPair.Modifiers = pair.Modifiers
			}
			pairs[key] = newPair
		}
	}

	return &object.Hash{Pairs: pairs}
}

func hasModifier(modifiers []int64, modifier int64) bool {
	for mod := range modifiers {
		if modifiers[mod] == modifier {
			return true
		}
	}
	return false
}

func MakeBuiltinClass(className string, fields []StringObjectPair) object.Hash {
	instance := MakeBuiltinInterface(fields)

	// instance.Constructor = instance.Pairs.Get(&object.String(className).HashKey())
	// instance.className = className
	// instance.Pairs.builtin = &object.HashPair{Key: strBuiltin.HashKey(), Value: TRUE}
	return *instance
}

type StringObjectPair struct {
	Name string
	Obj  object.Object
}

func MakeBuiltinInterface(methods []StringObjectPair) *object.Hash {
	pairs := make(map[object.HashKey]object.HashPair)
	for _, v := range methods {
		key := &object.String{Value: v.Name}
		pairs[key.HashKey()] = object.HashPair{
			Key:   key,
			Value: v.Obj,
		}
	}

	return &object.Hash{Pairs: pairs}
}

// func addMethod (allMethods, methodName string, contextName string, builtinFn object.Builtin) {
// 	allMethods = append(allMethods, &{[methodName]: &object.HashPair{
// 		Key: new object.String(methodName),
// 		Value: new object.Builtin(builtinFn, contextName)
// 	}});
// }

func NativeListToArray(items []interface{}) object.Array {
	var (
		elements []object.Object
	)
	for _, element := range items {
		switch element.(type) {
		case string:
			elements = append(elements, &object.String{Value: element.(string)})
		case int64:
			elements = append(elements, &object.Integer{Value: element.(int64)})
		case float64:
			elements = append(elements, &object.Float{Value: element.(float64)})
		case bool:
			elements = append(elements, &object.Boolean{Value: element.(bool)})
		// case []interface{}:
		// 	elements = append(elements, (nativeListToArray(element))
		// case interface{}:
		// 	elements = append(elements, nativeObjToMap(element.(map[string]interface{})).(interface{}(object.Object).(type)))
		default:
			elements = append(elements, &object.Null{})
		}
	}

	return object.Array{Elements: elements} //obj
}

// func nativeObjToMap (obj: {[key: string]: any} = {}): object.Hash => {
func NativeObjToMap(obj map[string]interface{}) object.Hash {
	newMap := object.Hash{Pairs: nil}

	for objectKey, data := range obj {
		var (
			value object.Object
		)

		switch data.(type) {
		case string:
			value = &object.String{Value: data.(string)}
			break
		case int64:
			value = &object.Integer{Value: data.(int64)}
			break
		case float64:
			value = &object.Float{Value: data.(float64)}
		case bool:
			value = &object.Boolean{Value: data.(bool)}
			break
		// case interface{}:
		// 	value = nativeObjToMap(data)
		// 	break
		// case :
		// 	console.log("native function", data)
		// 	// need to figure this out
		// 	// new object.Builtin(builtinFn, contextName)
		// 	break
		default:

		}
		key := &object.String{Value: objectKey}
		newMap.Pairs[key.HashKey()] = object.HashPair{
			Key:   key,
			Value: value,
		}
	}

	return newMap
}
