package evaluator

import (
	"math"

	"github.com/DeepAung/qcal/internal/object"
)

var builtinValues = map[string]object.Object{
	"pi": &object.Number{Value: math.Pi},
	"e":  &object.Number{Value: math.E},
}

type builtinFuncInfo struct {
	name  string
	len   int
	types []object.ObjectType
}

var infos = map[string]builtinFuncInfo{
	"min": {
		name:  "min",
		len:   2,
		types: []object.ObjectType{object.NUMBER_OBJ, object.NUMBER_OBJ},
	},
	"max": {
		name:  "max",
		len:   2,
		types: []object.ObjectType{object.NUMBER_OBJ, object.NUMBER_OBJ},
	},
	"abs":   {name: "abs", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"ceil":  {name: "ceil", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"floor": {name: "floor", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"round": {name: "round", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},

	"sqrt": {name: "sqrt", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"cbrt": {name: "cbrt", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},

	"log": {
		name:  "log",
		len:   2,
		types: []object.ObjectType{object.NUMBER_OBJ, object.NUMBER_OBJ},
	},
	"ln":    {name: "ln", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"log10": {name: "log10", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"log2":  {name: "log2", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"pow": {
		name:  "pow",
		len:   2,
		types: []object.ObjectType{object.NUMBER_OBJ, object.NUMBER_OBJ},
	},
	"pow10": {name: "pow10", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},

	"sin":     {name: "sin", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"cos":     {name: "cos", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"tan":     {name: "tan", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"sinh":    {name: "sinh", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"cosh":    {name: "cosh", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"tanh":    {name: "tanh", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"arcsin":  {name: "arcsin", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"arccos":  {name: "arccos", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"arctan":  {name: "arctan", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"arcsinh": {name: "arcsinh", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"arccosh": {name: "arccosh", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"arctanh": {name: "arctanh", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},

	"gamma": {name: "gamma", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	"hypot": {
		name:  "hypot",
		len:   2,
		types: []object.ObjectType{object.NUMBER_OBJ, object.NUMBER_OBJ},
	},

	// "Exp":         {name: "Exp", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	// "Exp2":        {name: "Exp2", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	// "Inf":         {name: "Inf", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	// "IsInf":       {name: "IsInf", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	// "IsNaN":       {name: "IsNaN", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
	// "Mod":         {name: "Mod", len: 2, types: []object.ObjectType{object.NUMBER_OBJ}},
	// "NaN":         {name: "NaN", len: 1, types: []object.ObjectType{object.NUMBER_OBJ}},
}

var builtinFuncs = map[string]object.BuiltinFunction{
	"min": func(args ...object.Object) object.Object {
		info := infos["min"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		val1 := args[1].(*object.Number).Value
		return newNumber(math.Min(val0, val1))
	},
	"max": func(args ...object.Object) object.Object {
		info := infos["max"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		val1 := args[1].(*object.Number).Value
		return newNumber(math.Max(val0, val1))
	},
	"abs": func(args ...object.Object) object.Object {
		info := infos["abs"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Abs(val0))
	},
	"ceil": func(args ...object.Object) object.Object {
		info := infos["ceil"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Ceil(val0))
	},
	"floor": func(args ...object.Object) object.Object {
		info := infos["floor"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Floor(val0))
	},
	"round": func(args ...object.Object) object.Object {
		info := infos["round"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Round(val0))
	},
	"sqrt": func(args ...object.Object) object.Object {
		info := infos["sqrt"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Sqrt(val0))
	},
	"cbrt": func(args ...object.Object) object.Object {
		info := infos["cbrt"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Cbrt(val0))
	},
	"log": func(args ...object.Object) object.Object {
		info := infos["log"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		val1 := args[1].(*object.Number).Value
		return newNumber(math.Log(val0) / math.Log(val1))
	},
	"ln": func(args ...object.Object) object.Object {
		info := infos["ln"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Log(val0))
	},
	"log10": func(args ...object.Object) object.Object {
		info := infos["log10"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Log10(val0))
	},
	"log2": func(args ...object.Object) object.Object {
		info := infos["log2"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Log2(val0))
	},
	"pow": func(args ...object.Object) object.Object {
		info := infos["pow"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		val1 := args[1].(*object.Number).Value
		return newNumber(math.Pow(val0, val1))
	},
	"pow10": func(args ...object.Object) object.Object {
		info := infos["pow10"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Pow(10, val0))
	},
	"sin": func(args ...object.Object) object.Object {
		info := infos["sin"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		_ = val0
		return newNumber(math.Sin(math.Pi))
	},
	"cos": func(args ...object.Object) object.Object {
		info := infos["cos"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Cos(val0))
	},
	"tan": func(args ...object.Object) object.Object {
		info := infos["tan"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Tan(val0))
	},
	"sinh": func(args ...object.Object) object.Object {
		info := infos["sinh"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Sinh(val0))
	},
	"cosh": func(args ...object.Object) object.Object {
		info := infos["cosh"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Cosh(val0))
	},
	"tanh": func(args ...object.Object) object.Object {
		info := infos["tanh"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Tanh(val0))
	},
	"arcsin": func(args ...object.Object) object.Object {
		info := infos["arcsin"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Asin(val0))
	},
	"arccos": func(args ...object.Object) object.Object {
		info := infos["arccos"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Acos(val0))
	},
	"arctan": func(args ...object.Object) object.Object {
		info := infos["arctan"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Atan(val0))
	},
	"arcsinh": func(args ...object.Object) object.Object {
		info := infos["arcsinh"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Asinh(val0))
	},
	"arccosh": func(args ...object.Object) object.Object {
		info := infos["arccosh"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Acosh(val0))
	},
	"arctanh": func(args ...object.Object) object.Object {
		info := infos["arctanh"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Atanh(val0))
	},
	"gamma": func(args ...object.Object) object.Object {
		info := infos["gamma"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		return newNumber(math.Gamma(val0))
	},
	"hypot": func(args ...object.Object) object.Object {
		info := infos["hypot"]
		if err := checkArgsLength(info, args); err != nil {
			return err
		}
		if err := checkArgsType(info, args); err != nil {
			return err
		}

		val0 := args[0].(*object.Number).Value
		val1 := args[1].(*object.Number).Value
		return newNumber(math.Hypot(val0, val1))
	},
}

func checkArgsLength(info builtinFuncInfo, args []object.Object) *object.Error {
	expect := info.len
	got := len(args)

	if got == -1 { // manually check args length
		return nil
	}

	if got < expect {
		return newError("%q: not enough arguments, expect=%d, got=%d", info.name, expect, got)
	} else if got > expect {
		return newError("%q: too many arguments, expect=%d, got=%d", info.name, expect, got)
	}
	return nil
}

func checkArgsType(info builtinFuncInfo, args []object.Object) *object.Error {
	for i, arg := range args {
		if arg.Type() != info.types[i] {
			return newError(
				"argument index %d of function %q should be type %s, got %s",
				i, info.name, info.types[i], arg.Type(),
			)
		}
	}
	return nil
}
