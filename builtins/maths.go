package builtins

import (
	"math"

	"github.com/SpaceHexagon/ecs/object"
	"github.com/SpaceHexagon/ecs/util"
)

func maths() *object.Hash {
	return util.MakeBuiltinInterface([]util.StringObjectPair{
		util.StringObjectPair{Name: "PI", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				return &object.Float{Value: math.Pi}
			},
		}},
		util.StringObjectPair{Name: "sin", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ {
					return newError("argument to `sin` must be FLOAT_OBJ, got %s", args[0].Type())
				}

				return &object.Float{Value: math.Sin(args[0].(*object.Float).Value)}
			},
		}},
		util.StringObjectPair{Name: "cos", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ {
					return newError("argument to `cos` must be FLOAT, got %s", args[0].Type())
				}
				return &object.Float{Value: math.Cos(args[0].(*object.Float).Value)}
			},
		}},
		util.StringObjectPair{Name: "tan", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ {
					return newError("argument to `cos` must be FLOAT, got %s", args[0].Type())
				}
				return &object.Float{Value: math.Tan(args[0].(*object.Float).Value)}
			},
		}},
		util.StringObjectPair{Name: "atan2", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ || args[1].Type() != object.FLOAT_OBJ {
					return newError("argument to `atan2` must be FLOAT, got %s %s", args[0].Type(), args[1].Type())
				}
				return &object.Float{Value: math.Atan2(args[0].(*object.Float).Value, args[1].(*object.Float).Value)}
			},
		}},
		util.StringObjectPair{Name: "sqrt", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ && args[0].Type() != object.INTEGER_OBJ {
					return newError("argument to `cos` must be INTEGER or FLOAT, got %s", args[0].Type())
				}
				return &object.Float{Value: math.Sqrt(float64(args[0].(*object.Integer).Value))}
			},
		}},
		util.StringObjectPair{Name: "abs", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ && args[0].Type() != object.INTEGER_OBJ {
					return newError("argument to `abs` must be INTEGER, got %s", args[0].Type())
				}
				return &object.Float{Value: math.Abs(float64(args[0].(*object.Integer).Value))}
			},
		}},
		util.StringObjectPair{Name: "floor", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ {
					return newError("argument to `abs` must be FLOAT_OBJ, got %s", args[0].Type())
				}
				return &object.Integer{Value: int64(math.Floor(args[0].(*object.Float).Value))}
			},
		}},
		util.StringObjectPair{Name: "ceil", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ {
					return newError("argument to `abs` must be FLOAT_OBJ, got %s", args[0].Type())
				}
				return &object.Float{Value: math.Ceil(args[0].(*object.Float).Value)}
			},
		}},
		util.StringObjectPair{Name: "fract", Obj: &object.Builtin{
			Fn: func(context interface{}, scope interface{}, args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != object.FLOAT_OBJ {
					return newError("argument to `abs` must be FLOAT_OBJ, got %s", args[0].Type())
				}
				return &object.Float{Value: args[0].(*object.Float).Value - math.Floor(args[0].(*object.Float).Value)}
			},
		}},
	})
}
