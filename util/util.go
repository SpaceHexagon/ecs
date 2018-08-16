package util

import (
	"github.com/SpaceHexagon/ecs/object"
)

func copyObject(valueNode object.Object) object.Object {
	switch valueNode.Type() {
	case "boolean":
		return &object.Boolean{Value: valueNode.(*object.Boolean).Value}
	case "int":
		return &object.Integer{Value: valueNode.(*object.Integer).Value}
	case "float":
		return &object.Float{Value: valueNode.(*object.Float).Value}
	case "string":
		return &object.String{Value: valueNode.(*object.String).Value}
	case "array":
		return &object.Array{Elements: valueNode.(*object.Array).Elements}
	case "hash":
		return copyHashMap(valueNode)
	case "function":
	case "BUILTIN":
		return valueNode
	default:
		return &object.Null{}
	}
	return &object.Null{}
}

func copyHashMap(data object.Object) object.Object {
	pairData := data.(*object.Hash).Pairs

	pairs := make(map[object.HashKey]object.HashPair)

	for key, pair := range pairData {
		valueNode := pair.Value
		keyNode := pair.Key
		isStatic := pair.Modifiers != nil && hasModifier(pair.Modifiers, 1)

		var (
			NewValue object.Object
			newPair  object.HashPair
		)

		if isStatic {
			// pairs.Set(KeyNode.HashKey(), valueNode)
		} else {
			NewValue = copyObject(valueNode)
			//newPair = {Key: keyNode, Value: NewValue} as object.HashPair;
			if pair.Modifiers != nil {
				newPair.Modifiers = pair.Modifiers
			}
			// pairs[(keyNode as any).Value] = newPair;
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

func makeBuiltinClass(className string, fields []StringObjectPair) object.Hash {
	instance := makeBuiltinInterface(fields)

	// instance.Constructor = instance.Pairs.Get(&object.String(className).HashKey())
	// instance.className = className
	// instance.Pairs.builtin = &object.HashPair{Key: strBuiltin.HashKey(), Value: TRUE}
	return instance
}

type StringObjectPair struct {
	name string
	obj  object.Object
}

func makeBuiltinInterface(methods []StringObjectPair) object.Hash {
	pairs := make(map[object.HashKey]object.HashPair)
	for k, v := range methods {
		v.Set(&object.HashKey, &object.HashPair{
			Key:   v.obj,
			Value: v,
		})
	}

	return &object.Hash{Pairs: pairs}
}

// func addMethod (allMethods, methodName string, contextName string, builtinFn object.Builtin) {
// 	allMethods = append(allMethods, &{[methodName]: &object.HashPair{
// 		Key: new object.String(methodName),
// 		Value: new object.Builtin(builtinFn, contextName)
// 	}});
// }

func nativeListToArray(obj slice) object.Array {
	return &object.Array{Elements: nil} //obj
	// 	.map(element => {
	// 		switch(typeof element) {
	// 			case "string":
	// 				return new object.String(element)
	// 			case "number":
	// 				return new object.Float(element)
	// 			case "boolean":
	// 				return new object.Boolean(element);
	// 			case "object":
	// 				if (typeof element.length == "number") {
	// 					return nativeListToArray(element);
	// 				} else {
	// 					return nativeObjToMap(element);
	// 				}
	// 			default:
	// 			return new object.NULL();
	// 		}
	// 	}
	// ));
}

// func nativeObjToMap (obj: {[key: string]: any} = {}): object.Hash => {
// func nativeObjToMap(obj map[string]interface) object.Hash {
// 	map := &object.Hash({ });

// 		for objectKey, data = range obj {
// 			var (
// 				value object.Object
// 			)

// 			switch(typeof data) {
// 				case "string":
// 					value = &object.String{Value: data}
// 				break;
// 				case "":
// 					value = &object.Integer{Value: data}
// 				break;
// 				case "boolean":
// 					value = &object.Boolean{Value: data}
// 				break;
// 				case "object":
// 					if (typeof data.length == "number") {
// 						value = nativeListToArray(data);
// 					} else {
// 						value = nativeObjToMap(data);
// 					}
// 				break;
// 				case "function":
// 				console.log("native function", data);
// 					// need to figure this out
// 					// new object.Builtin(builtinFn, contextName)
// 				break;
// 				default:

// 			}

// 			map[objectKey] = {
// 				Key: new object.String(objectKey),
// 				Value: value
// 			} as object.HashPair;
// 		}

// 		return map;
// 	}
