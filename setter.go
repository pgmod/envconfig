package envconfig

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func setValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		if value == "" {
			field.SetBool(false)
			return nil
		}
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(v)

	case reflect.Int, reflect.Int64:
		if value == "" {
			field.SetInt(0)
			return nil
		}
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(v)

	case reflect.Slice, reflect.Array:
		return setSliceOrArray(field, value)

	default:
		return fmt.Errorf("unsupported kind: %s", field.Kind())
	}

	return nil
}

func setSliceOrArray(field reflect.Value, value string) error {
	elemType := field.Type().Elem()

	// Проверяем, что элемент массива/слайса имеет тип int
	if elemType.Kind() != reflect.Int && elemType.Kind() != reflect.Int64 {
		return fmt.Errorf("unsupported slice/array element type: %s", elemType.Kind())
	}

	// Если значение пустое, создаем пустой слайс/массив
	if value == "" {
		if field.Kind() == reflect.Slice {
			field.Set(reflect.MakeSlice(field.Type(), 0, 0))
		} else {
			// Для массива оставляем нулевые значения
			for i := 0; i < field.Len(); i++ {
				field.Index(i).SetInt(0)
			}
		}
		return nil
	}

	// Разделяем строку по запятым
	parts := strings.Split(value, ",")

	// Для массива проверяем, что количество элементов совпадает
	if field.Kind() == reflect.Array {
		if len(parts) != field.Len() {
			return fmt.Errorf("array length mismatch: got %d values, expected %d", len(parts), field.Len())
		}
	}

	// Парсим каждое значение
	intValues := make([]int64, 0, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			intValues = append(intValues, 0)
			continue
		}
		v, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int value at index %d: %w", i, err)
		}
		intValues = append(intValues, v)
	}

	// Для слайса создаем новый слайс нужного размера
	if field.Kind() == reflect.Slice {
		slice := reflect.MakeSlice(field.Type(), len(intValues), len(intValues))
		for i, v := range intValues {
			slice.Index(i).SetInt(v)
		}
		field.Set(slice)
	} else {
		// Для массива устанавливаем значения
		for i, v := range intValues {
			field.Index(i).SetInt(v)
		}
	}

	return nil
}
