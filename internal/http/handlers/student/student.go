package student

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Abir-Zayn/studentsApi/internal/types"
	"github.com/Abir-Zayn/studentsApi/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var st types.Student

		// body is empty
		if r.Body == nil {
			slog.Error("empty body")
			response.SendBadRequest(w, "request body is empty")
			return
		}

		// decode json from req body
		err := json.NewDecoder(r.Body).Decode(&st)
		if err != nil {

			// Handling different types of error
			slog.Error("error decoding json", "error", err)

			switch err {
			case io.EOF:
				response.SendBadRequest(w, "request body is empty")
				return 
			case io.ErrUnexpectedEOF:
				response.SendBadRequest(w, "malformed json")
				return
			default:
				// if json has syntax error
				if _,ok := err.(*json.SyntaxError); ok {
					response.SendJSONParseError(w, fmt.Errorf("json syntax error: %w", err))
					return
				}
				// if json has type error
				if _,ok := err.(*json.UnmarshalTypeError); ok {
					response.SendBadRequest(w, "Invalid data type in JSON body")
					return
				}

				// other errors
				response.SendJSONParseError(w, fmt.Errorf("error parsing json: %w", err))
				return
			}
		}

		// validate student data required fields
		var missingFields []string
		if strings.TrimSpace(st.Name) == "" {
			missingFields = append(missingFields, "name")
		}
		if strings.TrimSpace(st.Email) == "" {
			missingFields = append(missingFields, "email")
		}

		// Send validation error if there are missing required fields
		if len(missingFields) > 0 {
			response.SendMissingFieldsError(w, missingFields)
			return
		}

		







	}
}