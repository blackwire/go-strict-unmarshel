package Example

import (
    "common"	
)

// Example is a test struct that explains how to use the method below
type Example struct {
    Field string `json:"field_name"`
    OmitableField string `json:"omitable_field,omitempty"`
}

func (e *Example) UnmarshalJSON(data []byte) error {
    return common.StrictUnmarshal(e, data)
}

// ---------------------------------------------------------------------------
// the below should be added to some common package that is imported
// ---------------------------------------------------------------------------

import (
    "encoding/json"
    "fmt"
    "reflect"
    "strings"
)

// StrictUnmarshal enforces unmarshal in a way that overwrites the original data should it exist
func StrictUnmarshal(s interface{}, data []byte) error {

	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	structValue := reflect.Indirect(reflect.ValueOf(s))
	structType := structValue.Type()

FieldLoop:
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := reflect.Indirect(reflect.ValueOf(s)).FieldByName(field.Name)

		for key, value := range m {
			if strings.ToLower(field.Name) == strings.ReplaceAll(key, "_", "") {
				if !reflect.TypeOf(value).ConvertibleTo(field.Type) {
					return fmt.Errorf("cannot convert field %s to %s", field.Name, field.Type)
				}

				fieldValue.Set(reflect.ValueOf(value))
				continue FieldLoop
			}
		}

		fieldValue.Set(reflect.Zero(fieldValue.Type()))
	}

	return nil
}
