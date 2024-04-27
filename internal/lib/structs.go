package lib

import (
	"fmt"
	"reflect"
)

func StructToMapSingleD(in interface{}, tag string) (map[string]interface{}, error) {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	out := make(map[string]interface{})
	queue := make([]interface{}, 0, 1)
	queue = append(queue, in)

	for len(queue) > 0 {
		v := reflect.ValueOf(queue[0])
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		queue = queue[1:]
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			vi := v.Field(i)
			if vi.Kind() == reflect.Ptr {
				vi = vi.Elem()
				if vi.Kind() == reflect.Struct {
					queue = append(queue, vi.Interface())
				} else {
					ti := t.Field(i)
					if tagValue := ti.Tag.Get(tag); tagValue != "" {
						out[tagValue] = vi.Interface()
					}
				}
				break
			}
			if vi.Kind() == reflect.Struct {
				queue = append(queue, vi.Interface())
				continue
			}
			ti := t.Field(i)
			if tagValue := ti.Tag.Get(tag); tagValue != "" {
				out[tagValue] = vi.Interface()
			}
		}
	}
	return out, nil
}

func StructToMapSingleD2(in interface{}, tag string, result *map[string]interface{}) error {
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	queue := make([]interface{}, 0, 1)
	queue = append(queue, in)

	for len(queue) > 0 {
		v := reflect.ValueOf(queue[0])
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		queue = queue[1:]
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			vi := v.Field(i)
			if vi.Kind() == reflect.Ptr {
				vi = vi.Elem()
				if vi.Kind() == reflect.Struct {
					queue = append(queue, vi.Interface())
				} else {
					ti := t.Field(i)
					if tagValue := ti.Tag.Get(tag); tagValue != "" {
						(*result)[tagValue] = vi.Interface()
					}
				}
				break
			}
			if vi.Kind() == reflect.Struct {
				queue = append(queue, vi.Interface())
				continue
			}
			ti := t.Field(i)
			if tagValue := ti.Tag.Get(tag); tagValue != "" {
				(*result)[tagValue] = vi.Interface()
			}
		}
	}
	return nil
}
