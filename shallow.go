package clone

import (
	"reflect"
)

func Shallow(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}

	src := reflect.ValueOf(obj)
	return shallow(src).Interface()
}

func shallow(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Invalid:
		return reflect.ValueOf(nil)

	case reflect.Ptr:
		if v.IsNil() {
			res := reflect.New(v.Type())
			return res.Elem()
		}

		res := reflect.New(v.Elem().Type())
		res.Elem().Set(v.Elem())
		return res

	case reflect.Map:
		res := reflect.MakeMap(
			reflect.MapOf(
				v.Type().Key(),
				v.Type().Elem()))

		for _, k := range v.MapKeys() {
			value := v.MapIndex(k)
			res.SetMapIndex(k, value)
		}

		return res

	case reflect.Array:
		len := v.Len()
		types := reflect.ArrayOf(len, v.Type().Elem())
		res := reflect.New(types).Elem()

		for i := 0; i < len; i++ {
			res.Index(i).Set(v.Index(i))
		}

		return res

	case reflect.Slice:
		len := v.Len()
		cap := v.Cap()
		res := reflect.MakeSlice(v.Type(), len, cap)

		for i := 0; i < len; i++ {
			res.Index(i).Set(v.Index(i))
		}

		return res

	default:
		res := reflect.New(v.Type()).Elem()
		res.Set(v)
		return res
	}
}
