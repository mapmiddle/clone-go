package clone

import (
	"reflect"
	"unsafe"
)

func Deep(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}

	v := reflect.ValueOf(obj)
	return deep(v).Interface()
}

func deep(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		return deepPtr(v)
	case reflect.Struct:
		return deepStruct(v)
	case reflect.Map:
		return deepMap(v)
	case reflect.Array:
		return deepArray(v)
	case reflect.Slice:
		return deepSlice(v)
	// primitive
	default:
		if v.CanInterface() {
			res := reflect.New(v.Type()).Elem()
			res.Set(v)
			return res
		}

		if v.CanAddr() {
			res := reflect.New(v.Type()).Elem()
			res.Set(unexportedValue(v))
			return res
		}

		res := v
		return res
	}
}

func deepMap(v reflect.Value) reflect.Value {
	res := reflect.MakeMap(
		reflect.MapOf(
			v.Type().Key(),
			v.Type().Elem()))

	for _, k := range v.MapKeys() {
		key := k.Convert(res.Type().Key())
		value := v.MapIndex(key)
		elem := deep(value.Elem())
		res.SetMapIndex(key, elem)
	}

	return res
}

func deepArray(v reflect.Value) reflect.Value {
	len := v.Len()
	t := reflect.ArrayOf(len, v.Type().Elem())
	res := reflect.New(t).Elem()

	for i := 0; i < len; i++ {
		dest := res.Index(i)
		src := deep(v.Index(i))
		dest.Set(src)
	}

	return res
}

func deepSlice(v reflect.Value) reflect.Value {
	len := v.Len()
	cap := v.Cap()
	res := reflect.MakeSlice(v.Type(), len, cap)

	for i := 0; i < len; i++ {
		dest := res.Index(i)
		src := deep(v.Index(i))
		dest.Set(src)
	}

	return res
}

func deepPtr(v reflect.Value) reflect.Value {
	if v.IsNil() {
		res := reflect.New(v.Type())
		return res.Elem()
	}

	elem := deep(v.Elem())
	return elem.Addr()
}

func deepStruct(v reflect.Value) reflect.Value {
	ptr := reflect.New(v.Type())
	elem := ptr.Elem()

	if v.CanInterface() {
		elem.Set(v)
	}

	for i := 0; i < v.NumField(); i++ {
		src := deep(v.Field(i))
		dest := elem.Field(i)

		switch {
		case dest.CanInterface():
			dest.Set(src)
		case src.CanAddr():
			setUnexportedValue(dest, src)
		case src.Kind() == reflect.Ptr:
			setUnexportedValue(dest, src)
		}
	}

	return elem
}

func unexportedValue(f reflect.Value) reflect.Value {
	return reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).
		Elem()
}

func setUnexportedValue(f reflect.Value, v reflect.Value) {
	reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).
		Elem().
		Set(v)
}
