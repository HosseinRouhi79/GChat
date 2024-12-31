package helper

import "encoding/json"

func TypeConverter[T any](data any) (model *T, err error) {

	dataJson, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dataJson, &model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
