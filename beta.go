package itertools

import (
	"reflect"
)

// 支持Array,Slice
func Map(fn interface{}, sequences ...interface{}) []interface{} {
	mustOverZero(sequences...)
	mustBeFunc(fn)
	mustBeSeq(sequences...)
	mustBeEquilong(sequences...)
	mustMatchFunc(fn, 1, sequences...)

	argNum := len(sequences)
	length := reflect.ValueOf(sequences[0]).Len()
	fnValue := reflect.ValueOf(fn)
	res := make([]interface{}, length)
	for i := 0; i < length; i++ {
		args := make([]reflect.Value, 0, argNum)
		for _, sequence := range sequences {
			seqValue := reflect.ValueOf(sequence)
			args = append(args, indexOf(seqValue, i))
		}
		res[i] = fnValue.Call(args)[0].Interface()
	}
	return res
}

func Filter(fn interface{}, sequence interface{}) []interface{} {
	mustBeFunc(fn)
	mustBeSeq(sequence)
	mustMatchFunc(fn, 1, sequence)

	seqValue := reflect.ValueOf(sequence)
	length := seqValue.Len()
	fnValue := reflect.ValueOf(fn)
	res := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		args := make([]reflect.Value, 1)
		args[0] = indexOf(seqValue, i)
		if fnValue.Call(args)[0].Bool() {
			res = append(res, args[0].Interface())
		}
	}
	return res
}

func Reduce(fn interface{}, sequence interface{}, initial interface{}) interface{} {
	mustBeFunc(fn)
	mustBeSeq(sequence)
	mustBeFolded(sequence, initial)
	mustMatchFuncInit(fn, sequence, initial)

	seqValue := reflect.ValueOf(sequence)
	fnValue := reflect.ValueOf(fn)
	var res reflect.Value
	var start int
	if initial == nil {
		res = indexOf(seqValue, 0)
		start = 1
	} else {
		res = reflect.ValueOf(initial)
		start = 0
	}

	for i := start; i < seqValue.Len(); i++ {
		args := make([]reflect.Value, 2)
		args[0] = res
		args[1] = indexOf(seqValue, i)
		res = fnValue.Call(args)[0]
	}
	return res.Interface()
}

func Foreach(fn interface{}, sequence interface{}) {
	mustBeFunc(fn)
	mustBeSeq(sequence)
	mustMatchFunc(fn, 0, sequence)

	seqValue := reflect.ValueOf(sequence)
	length := seqValue.Len()
	fnValue := reflect.ValueOf(fn)
	for i := 0; i < length; i++ {
		args := make([]reflect.Value, 1)
		args[0] = indexOf(seqValue, i)
		fnValue.Call(args)
	}
}

func mustOverZero(objs ...interface{}) {
	if len(objs) == 0 {
		panic("length is 0")
	}
}

func mustBeSeq(objs ...interface{}) {
	for _, obj := range objs {
		mustBe(obj, reflect.Array, reflect.Slice, reflect.Chan)
	}
}

func mustBeFunc(objs ...interface{}) {
	for _, obj := range objs {
		mustBe(obj, reflect.Func)
	}
}

func mustBe(obj interface{}, expects ...reflect.Kind) {
	kind := reflect.ValueOf(obj).Kind()
	for _, expect := range expects {
		if kind == expect {
			return
		}
	}
	panic("unexpected kind")
}

func mustBeEquilong(objs ...interface{}) {
	for len(objs) == 0 {
		return
	}

	length := reflect.ValueOf(objs[0]).Len()
	for i := 1; i < len(objs); i++ {
		if length != reflect.ValueOf(objs[i]).Len() {
			panic("length is not equal")
		}
	}
}

func mustMatchFunc(fn interface{}, fnNumOut int, objs ...interface{}) {
	inKinds := make([]reflect.Kind, len(objs))
	for index, obj := range objs {
		inKinds[index] = indexOf(reflect.ValueOf(obj), 0).Kind()
	}

	fnType := reflect.TypeOf(fn)
	numIn := fnType.NumIn()
	numOut := fnType.NumOut()
	if numIn != len(inKinds) || numOut != fnNumOut {
		panic("function does not match")
	}

	for i := 0; i < numIn; i++ {
		if fnType.In(i).Kind() != inKinds[i] {
			panic("function does not match")
		}
	}
}

func mustBeFolded(obj interface{}, initial interface{}) {
	if initial == nil && reflect.ValueOf(obj).Len() == 0 {
		panic("empty sequence with no initial value")
	}
}

func mustMatchFuncInit(fn interface{}, obj interface{}, initial interface{}) {
	fnType := reflect.TypeOf(fn)
	numIn := fnType.NumIn()
	numOut := fnType.NumOut()
	if numIn != 2 || numOut != 1 {
		panic("function does not match")
	}

	fnIn1 := fnType.In(0).Kind()
	fnIn2 := fnType.In(1).Kind()
	fnOut := fnType.Out(0).Kind()
	objKind := indexOf(reflect.ValueOf(obj), 0).Kind()
	if initial == nil {
		if fnIn1 != objKind || fnIn2 != objKind || fnOut != objKind {
			panic("function does not match")
		}
	} else {
		initialKind := reflect.ValueOf(initial).Kind()
		if fnIn1 != fnOut || fnIn1 != initialKind || fnIn2 != objKind {
			panic("function does not match")
		}
	}
}

func indexOf(v reflect.Value, i int) reflect.Value {
	return getRealValue(v.Index(i))
}

func getRealValue(v reflect.Value) reflect.Value {
	return reflect.ValueOf(v.Interface())
}
