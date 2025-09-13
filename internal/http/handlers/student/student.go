package student

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/Abir-Zayn/studentsApi/internal/storage"
	"github.com/Abir-Zayn/studentsApi/internal/types"
	"github.com/Abir-Zayn/studentsApi/internal/utils/response"
	"github.com/google/uuid"
)

func New(storage storage.Storage) http.HandlerFunc {
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
		// Validate age 
		if st.Age < 0 || st.Age >= 20 {
			response.SendInvalidFieldError(w, "age", "age must be between 0 and 20")
			return
		}
		// Validate email format
		if !isValidEmail(st.Email) {
			response.SendInvalidFieldError(w, "email", "invalid email format")
			return
		}
		//validate group
		if st.Group != "" && !st.Group.IsValid() {
			response.SendInvalidFieldError(w, "group", "must be one of: Science, Arts, Commerce")
			return
		}

		// Generate student ID
		st.Id = generateStudentID()

		// Save student to database
		savedID, err := storage.CreateStudent(&st)
		if err != nil {
			slog.Error("failed to save student to database", "error", err)
			response.SendInternalServerError(w, "Failed to save student")
			return
		}

		//creating successful response
		studentData := map[string]interface{}{
			"student" : st,
			"id" : savedID,
		}
		response.SendCreated(w, "Student Added Successfully", studentData)
	}	
}



// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// generateStudentID generates a unique ID for the student
func generateStudentID() string {
	id := uuid.New()
	return id.String()
}