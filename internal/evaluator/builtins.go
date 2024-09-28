package evaluator

import (
	"math"

	"github.com/DeepAung/qcal/internal/object"
)

var builtinValues = map[string]object.Object{
	"pi": &object.Number{Value: math.Pi},
	"e":  &object.Number{Value: math.E},
}

var builtinFuncs = map[string]*object.BuiltinFunction{
	// "len": {
	// 	Fn: func(args ...object.Object) object.Object {
	// 		if len(args) != 1 {
	// 			return newError("wrong number of arguments. got=%d, want=1", len(args))
	// 		}
	//
	// 		switch arg := args[0].(type) {
	// 		case *object.String:
	// 			return &object.Integer{Value: int64(len(arg.Value))}
	// 		case *object.Array:
	// 			return &object.Integer{Value: int64(len(arg.Elements))}
	// 		default:
	// 			return newError("argument to `len` not supported, got %s", args[0].Type())
	// 		}
	// 	},
	// },
	// "first": {
	// 	Fn: func(args ...object.Object) object.Object {
	// 		if len(args) != 1 {
	// 			return newError("wrong number of arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `first` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}
	//
	// 		arr := args[0].(*object.Array)
	// 		if len(arr.Elements) == 0 {
	// 			return NULL
	// 		}
	//
	// 		return arr.Elements[0]
	// 	},
	// },
	// "last": {
	// 	Fn: func(args ...object.Object) object.Object {
	// 		if len(args) != 1 {
	// 			return newError("wrong number of arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `last` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}
	//
	// 		arr := args[0].(*object.Array)
	// 		length := len(arr.Elements)
	// 		if length == 0 {
	// 			return NULL
	// 		}
	//
	// 		return arr.Elements[length-1]
	// 	},
	// },
	// "rest": {
	// 	Fn: func(args ...object.Object) object.Object {
	// 		if len(args) != 1 {
	// 			return newError("wrong number of arguments. got=%d, want=1",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `rest` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}
	//
	// 		arr := args[0].(*object.Array)
	// 		length := len(arr.Elements)
	// 		if length == 0 {
	// 			return NULL
	// 		}
	//
	// 		newElements := make([]object.Object, length-1, length-1)
	// 		copy(newElements, arr.Elements[1:length])
	//
	// 		return &object.Array{Elements: newElements}
	// 	},
	// },
	// "push": {
	// 	Fn: func(args ...object.Object) object.Object {
	// 		if len(args) != 2 {
	// 			return newError("wrong number of arguments. got=%d, want=2",
	// 				len(args))
	// 		}
	// 		if args[0].Type() != object.ARRAY_OBJ {
	// 			return newError("argument to `push` must be ARRAY, got %s",
	// 				args[0].Type())
	// 		}
	//
	// 		arr := args[0].(*object.Array)
	// 		length := len(arr.Elements)
	//
	// 		newElements := make([]object.Object, length+1, length+1)
	// 		copy(newElements, arr.Elements)
	// 		newElements[length] = args[1]
	//
	// 		return &object.Array{Elements: newElements}
	// 	},
	// },
}
