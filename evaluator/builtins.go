package evaluator

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/SpaceHexagon/ecs/object"
)

var builtins = map[string]*object.Builtin{
	"PI": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			return &object.Float{Value: math.Pi}
		},
	},
	"sin": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ {
				return newError("argument to `sin` must be INTEGER, got %s", args[0].Type())
			}

			return &object.Float{Value: math.Sin(args[0].(*object.Float).Value)}
		},
	},
	"cos": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ {
				return newError("argument to `cos` must be FLOAT_OBJ, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Cos(args[0].(*object.Float).Value)}
		},
	},
	"tan": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ {
				return newError("argument to `cos` must be FLOAT_OBJ, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Tan(args[0].(*object.Float).Value)}
		},
	},
	"atan2": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ || args[1].Type() != object.FLOAT_OBJ {
				return newError("argument to `atan2` must be FLOAT_OBJ, got %s %s", args[0].Type(), args[1].Type())
			}
			return &object.Float{Value: math.Atan2(args[0].(*object.Float).Value, args[1].(*object.Float).Value)}
		},
	},
	"sqrt": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ && args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `cos` must be INTEGER or FLOAT, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sqrt(float64(args[0].(*object.Integer).Value))}
		},
	},
	"abs": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ && args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `abs` must be INTEGER, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Abs(float64(args[0].(*object.Integer).Value))}
		},
	},
	"floor": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ {
				return newError("argument to `abs` must be FLOAT_OBJ, got %s", args[0].Type())
			}
			return &object.Float{Value: args[0].(*object.Float).Value - math.Floor(args[0].(*object.Float).Value)}
		},
	},
	"ceil": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ {
				return newError("argument to `abs` must be FLOAT_OBJ, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Ceil(args[0].(*object.Float).Value)}
		},
	},
	"fract": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.FLOAT_OBJ {
				return newError("argument to `abs` must be FLOAT_OBJ, got %s", args[0].Type())
			}
			return &object.Float{Value: args[0].(*object.Float).Value - math.Floor(args[0].(*object.Float).Value)}
		},
	},
	"time": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {

			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"print": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"len": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}
			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
	"join": &object.Builtin{
		Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("first argument to `join` must be ARRAY, got %s",
					args[0].Type())
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("second argument to `join` must be STRING, got %s",
					args[0].Type())
			}
			strArray := []string{}
			arr := args[0].(*object.Array)
			for _, element := range arr.Elements {
				s := ""
				s = element.Inspect()
				strArray = append(strArray, s)
			}
			outStr := strings.Join(strArray, args[1].(*object.String).Value)
			return &object.String{Value: outStr}
		},
	},
}
