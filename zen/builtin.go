package zen

func builtinNth(v *Value, n int) *Value {
	for ; n > 0; n-- {
		v = v.Pair().cdr
	}
	return v.Pair().car
}

func builtinCons(a, b *Value) *Value {
	return NewPairValue(a, b)
}

func compareValue(v1, v2 *Value) int {
	if v1 == v2 {
		return 0
	}
	if v1.Type() == ValueTypeInteger && v2.Type() == ValueTypeInteger {
		return v1.Int() - v2.Int()
	}
	panic("not support yet")
}

func compare(v1, v2 *Value) *Value {
	return NewIntegerValue(compareValue(v1, v2))
}

func equal(v1, v2 *Value) *Value {
	if compareValue(v1, v2) == 0 {
		return True
	}
	return False
}

func notequal(v1, v2 *Value) *Value {
	if compareValue(v1, v2) == 0 {
		return False
	}
	return True
}

func less(v1, v2 *Value) *Value {
	if compareValue(v1, v2) < 0 {
		return True
	}
	return False
}

func lessequal(v1, v2 *Value) *Value {
	if compareValue(v1, v2) <= 0 {
		return True
	}
	return False
}

func greater(v1, v2 *Value) *Value {
	if compareValue(v1, v2) > 0 {
		return True
	}
	return False
}

var cprims = []func(v1, v2 *Value) *Value{
	compare,
	equal,
	notequal,
	less,
	lessequal,
	greater,
}
