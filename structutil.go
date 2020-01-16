package expressgo

import (
	"fmt"
	"reflect"
)

func setStructFromMap(st interface{}, src map[string]interface{}) {
	structValue := reflect.ValueOf(st).Elem()

	setField := func(obj interface{}, name string, value interface{}) error {

		structFieldValue := structValue.FieldByName(name)

		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in %s", name, structValue.Type().Name())
		}

		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value in %s", name, structValue.Type().Name())
		}

		structFieldType := structFieldValue.Type()
		val := reflect.ValueOf(value)
		if structFieldType != val.Type() {
			return fmt.Errorf("Invalid type for %s field value in %s. Expected %s", name, structValue.Type().Name(), structFieldValue.Type().Name())
		}

		structFieldValue.Set(val)
		return nil
	}

	for k, v := range src {
		if e := setField(st, k, v); e != nil {
			panic(e)
		}
	}
}
