package response

import "net/http"
func Success(message string, data interface{}) Response {
	return Response{
		Status:     StatusSuccess,
		Message:    message,
		Data:       data,
		StatusCode: http.StatusOK,
	}
}

func Created(message string, data interface{}) Response {
	return Response{
		Status:     StatusSuccess,
		Message:    message,
		Data:       data,
		StatusCode: http.StatusCreated,
	}
}

func Accepted(message string, data interface{}) Response {
	return Response{
		Status:     StatusSuccess,
		Message:    message,
		Data:       data,
		StatusCode: http.StatusAccepted,
	}
}

func NoContent(message string) Response {
	return Response{
		Status:     StatusSuccess,
		Message:    message,
		StatusCode: http.StatusNoContent,
	}
}

// Quick send methods for success responses
func SendSuccess(w http.ResponseWriter, message string, data interface{}) error {
	return Send(w, Success(message, data))
}

func SendCreated(w http.ResponseWriter, message string, data interface{}) error {
	return Send(w, Created(message, data))
}

func SendAccepted(w http.ResponseWriter, message string, data interface{}) error {
	return Send(w, Accepted(message, data))
}

func SendNoContent(w http.ResponseWriter, message string) error {
	return Send(w, NoContent(message))
}
