package validation

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/Abir-Zayn/studentsApi/internal/types"
)

// For specific validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Return validation error results
type ValidationResults struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors"`
}

// Add a validation error to the results
func (vr* ValidationResults) Add(field, message string) {
	vr.IsValid = false
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}


// GetFieldErrors returns all error messages for a specific field
func (vr *ValidationResults) GetFieldErrors(field string) []string {
	var errors []string
	for _, err := range vr.Errors {
		if err.Field == field {
			errors = append(errors, err.Message)
		}
	}
	return errors
}

// ValidateStudent performs comprehensive validation on a Student struct
func ValidateStudent(student types.Student) ValidationResults {
	result := ValidationResults{IsValid: true}

	// Validate required fields
	if strings.TrimSpace(student.Name) == "" {
		result.Add("name", "name is required")
	} else if len(student.Name) > 100 {
		result.Add("name", "name must not exceed 100 characters")
	}

	if strings.TrimSpace(student.Email) == "" {
		result.Add("email", "email is required")
	} else if !IsValidEmail(student.Email) {
		result.Add("email", "invalid email format")
	}

	// Validate optional fields
	if student.Age != 0 && (student.Age < 5 || student.Age > 100) {
		result.Add("age", "age must be between 5 and 100")
	}

	if student.Group != "" && !student.Group.IsValid() {
		result.Add("group", "group must be one of: Science, Arts, Commerce")
	}

	if student.TutionFee < 0 {
		result.Add("tution_fee", "tuition fee cannot be negative")
	}

	if student.Phone != "" && !IsValidPhone(student.Phone) {
		result.Add("phone", "invalid phone number format")
	}

	if len(student.Address) > 500 {
		result.Add("address", "address must not exceed 500 characters")
	}

	if len(student.Mentor) > 100 {
		result.Add("mentor", "mentor name must not exceed 100 characters")
	}

	// Validate subjects array
	if len(student.Subjects) > 10 {
		result.Add("subjects", "cannot have more than 10 subjects")
	}

	for i, subject := range student.Subjects {
		if strings.TrimSpace(subject) == "" {
			result.Add("subjects", "subject at index "+string(rune(i))+" cannot be empty")
		}
		if len(subject) > 50 {
			result.Add("subjects", "subject at index "+string(rune(i))+" must not exceed 50 characters")
		}
	}

	return result
}

// IsValidEmail validates email format using Go's mail package
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsValidPhone validates phone number format (basic validation)
func IsValidPhone(phone string) bool {
	// Remove common formatting characters
	cleaned := regexp.MustCompile(`[\s\-\(\)\+]`).ReplaceAllString(phone, "")
	
	// Check if it contains only digits and is of reasonable length
	phoneRegex := regexp.MustCompile(`^\d{10,15}$`)
	return phoneRegex.MatchString(cleaned)
}

// ValidateRequiredFields checks if required fields are present
func ValidateRequiredFields(student types.Student) []string {
	var missingFields []string

	if strings.TrimSpace(student.Name) == "" {
		missingFields = append(missingFields, "name")
	}

	if strings.TrimSpace(student.Email) == "" {
		missingFields = append(missingFields, "email")
	}

	return missingFields
}

// ValidateUpdateFields validates fields for update operations (where not all fields are required)
func ValidateUpdateFields(student types.Student) ValidationResults {
	result := ValidationResults{IsValid: true}

	// Only validate fields that are provided (non-empty/non-zero)
	if student.Name != "" && len(student.Name) > 100 {
		result.Add("name", "name must not exceed 100 characters")
	}

	if student.Email != "" && !IsValidEmail(student.Email) {
		result.Add("email", "invalid email format")
	}

	if student.Age != 0 && (student.Age < 5 || student.Age > 20) {
		result.Add("age", "age must be between 5 and 20")
	}

	if student.Group != "" && !student.Group.IsValid() {
		result.Add("group", "group must be one of: Science, Arts, Commerce")
	}

	if student.TutionFee < 0 && student.TutionFee != 0 && student.TutionFee > 100000 {
		result.Add("tution_fee", "tuition fee cannot be negative or exceed 100000")
	}

	if student.Phone != "" && !IsValidPhone(student.Phone) {
		result.Add("phone", "invalid phone number format")
	}

	if len(student.Address) > 500 {
		result.Add("address", "address must not exceed 500 characters")
	}

	if len(student.Mentor) > 100 {
		result.Add("mentor", "mentor name must not exceed 100 characters")
	}

	if len(student.Subjects) > 10 {
		result.Add("subjects", "cannot have more than 10 subjects")
	}

	return result
}

