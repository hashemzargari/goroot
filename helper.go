package goroot

import (
	"os"
	"reflect"
)

// Write the content to a file.
func writeToFile(fileName, content string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

// getTypeName returns the name of the type.
func getTypeName(i interface{}) string {
	reflectType := reflect.TypeOf(i)
	if reflectType.Kind() == reflect.Array || reflectType.Kind() == reflect.Slice {
		return getTypeName(reflect.TypeOf(i).Elem())
	}
	if reflectType.Kind() == reflect.Map {
		return getTypeName(reflect.TypeOf(i).Elem())
	}
	return reflectType.Name()
}

// getFieldNamesAndTypes returns the name and type of the fields of the struct.
func getFieldNamesAndTypes(i interface{}) []struct {
	Name  string
	Type  reflect.Type
	Index int
} {
	var fields []struct {
		Name  string
		Type  reflect.Type
		Index int
	}
	t := reflect.TypeOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fields = append(fields, struct {
			Name  string
			Type  reflect.Type
			Index int
		}{
			field.Name,
			field.Type,
			i + 1,
		})
	}

	return fields
}
