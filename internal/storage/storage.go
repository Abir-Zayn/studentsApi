package storage

import "github.com/Abir-Zayn/studentsApi/internal/types"

type Storage interface {
	  CreateStudent(student *types.Student) (string, error)
	  GetStudentByID(id string) (*types.Student, error)
	  UpdateStudent(student *types.Student) error
	  DeleteStudent(id string) error
	  ListStudents() ([]*types.Student, error)
}