package helper

import "chat/src/api/validation"

type HTTPResponse struct {
	Result           any                          `json:"result"`
	Success          bool                         `json:"success"`
	ResultCode       string                       `json:"resultCode"`
	ValidationErrors []validation.ValidationError `json:"validationErrors"`
	Error            any                          `json:"error"`
}

func Response(result any, resultCode string) *HTTPResponse {
	var res HTTPResponse
	res.Result = result
	res.ResultCode = resultCode
	res.Success = true

	return &res

}
func ResponseWithError(result any, resultCode string, err error) *HTTPResponse {
	var res HTTPResponse
	res.Result = result
	res.ResultCode = resultCode
	res.Success = false
	res.Error = err.Error()

	return &res
}
func ResponseWithValidationError(result any, resultCode string, err error, veArr []validation.ValidationError) *HTTPResponse {
	arr := validation.GetValidationError(err)
	var res HTTPResponse
	res.Result = result
	res.ResultCode = resultCode
	res.Success = false
	res.ValidationErrors = *arr

	return &res
}
