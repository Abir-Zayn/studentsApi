package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Abir-Zayn/studentsApi/internal/config"
	"github.com/Abir-Zayn/studentsApi/internal/types"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New (cfg config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS students (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		age INTEGER DEFAULT 0,
		city TEXT DEFAULT '',
		email TEXT NOT NULL UNIQUE,
		"group" TEXT DEFAULT '',   -- Quoted because 'group' is a reserved keyword
		phone TEXT DEFAULT '',
		address TEXT DEFAULT '',
		tution_fee REAL DEFAULT 0.0,
		enrolled INTEGER DEFAULT 0,
		mentor TEXT DEFAULT '',
		subjects TEXT DEFAULT '[]', -- JSON array of strings
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		return nil, err
	}
	

	//TODO: create indexing
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_students_email ON students(email)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create email index: %w", err)
	}


	

	return  &Sqlite{
		Db: db,
	}, nil
}

// CreateStudent implements the Storage interface
func (s *Sqlite) CreateStudent(student *types.Student) (string, error) {
	// Convert subjects slice to JSON string
	subjectsJSON, err := json.Marshal(student.Subjects)
	if err != nil {
		return "", fmt.Errorf("failed to marshal subjects: %w", err)
	}

	query := `INSERT INTO students (id, name, age, city, email, "group", phone, address, tution_fee, enrolled, mentor, subjects) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err = s.Db.Exec(query, 
		student.Id, student.Name, student.Age, student.City, student.Email, 
		string(student.Group), student.Phone, student.Address, student.TutionFee, 
		student.Enrolled, student.Mentor, string(subjectsJSON))
	
	if err != nil {
		return "", fmt.Errorf("failed to insert student: %w", err)
	}

	return student.Id, nil
}

// GetStudentByID implements the Storage interface
func (s *Sqlite) GetStudentByID(id string) (*types.Student, error) {
	query := `SELECT id, name, age, city, email, "group", phone, address, tution_fee, enrolled, mentor, subjects 
              FROM students WHERE id = ?`
	
	var student types.Student
	var subjectsJSON string
	var group string
	
	err := s.Db.QueryRow(query, id).Scan(
		&student.Id, &student.Name, &student.Age, &student.City, &student.Email,
		&group, &student.Phone, &student.Address, &student.TutionFee,
		&student.Enrolled, &student.Mentor, &subjectsJSON)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student with id %s not found", id)
		}
		return nil, fmt.Errorf("failed to get student: %w", err)
	}

	student.Group = types.GroupType(group)
	
	// Parse subjects JSON
	if subjectsJSON != "" {
		err = json.Unmarshal([]byte(subjectsJSON), &student.Subjects)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal subjects: %w", err)
		}
	}

	return &student, nil
}

// UpdateStudent implements the Storage interface
func (s *Sqlite) UpdateStudent(student *types.Student) error {
	subjectsJSON, err := json.Marshal(student.Subjects)
	if err != nil {
		return fmt.Errorf("failed to marshal subjects: %w", err)
	}

	query := `UPDATE students SET name=?, age=?, city=?, email=?, "group"=?, phone=?, address=?, 
              tution_fee=?, enrolled=?, mentor=?, subjects=?, updated_at=CURRENT_TIMESTAMP 
              WHERE id=?`
	
	result, err := s.Db.Exec(query,
		student.Name, student.Age, student.City, student.Email, string(student.Group),
		student.Phone, student.Address, student.TutionFee, student.Enrolled,
		student.Mentor, string(subjectsJSON), student.Id)
	
	if err != nil {
		return fmt.Errorf("failed to update student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student with id %s not found", student.Id)
	}

	return nil
}

// DeleteStudent implements the Storage interface
func (s *Sqlite) DeleteStudent(id string) error {
	query := `DELETE FROM students WHERE id = ?`
	
	result, err := s.Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete student: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("student with id %s not found", id)
	}

	return nil
}

// ListStudents implements the Storage interface
func (s *Sqlite) ListStudents() ([]*types.Student, error) {
	query := `SELECT id, name, age, city, email, "group", phone, address, tution_fee, enrolled, mentor, subjects 
              FROM students ORDER BY created_at DESC`
	
	rows, err := s.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list students: %w", err)
	}
	defer rows.Close()

	var students []*types.Student
	
	for rows.Next() {
		var student types.Student
		var subjectsJSON string
		var group string
		
		err := rows.Scan(
			&student.Id, &student.Name, &student.Age, &student.City, &student.Email,
			&group, &student.Phone, &student.Address, &student.TutionFee,
			&student.Enrolled, &student.Mentor, &subjectsJSON)
		
		if err != nil {
			return nil, fmt.Errorf("failed to scan student: %w", err)
		}

		student.Group = types.GroupType(group)
		
		// Parse subjects JSON
		if subjectsJSON != "" {
			err = json.Unmarshal([]byte(subjectsJSON), &student.Subjects)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal subjects: %w", err)
			}
		}

		students = append(students, &student)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return students, nil
}