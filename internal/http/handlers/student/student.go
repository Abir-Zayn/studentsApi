package student

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Abir-Zayn/studentsApi/internal/types"
	"github.com/Abir-Zayn/studentsApi/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var st types.Student

		// reject empty body
		if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
			// io.EOF, io.UnexpectedEOF, syntax errors … all land here
			response.WriteJson(w, http.StatusBadRequest,
				response.GeneralError(
					fmt.Errorf("empty Body: %w", err)))
			return
		}

		//validate required fields
		if st.Name == "" || st.Email == "" {
			response.WriteJson(w, http.StatusBadRequest,
				map[string]string{"error": "missing required fields"})
			return
		}

		slog.Info("creating new student", "name", st.Name)
		response.WriteJson(w, http.StatusCreated,
			map[string]string{"success": "OK"})
	}
}