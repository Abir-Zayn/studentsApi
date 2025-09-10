package types

//Defining a custom type for GroupType
type GroupType string 

// Defining allowed constants
const  (
	GroupScience GroupType = "Science"
	GroupArts   GroupType = "Arts"
	GroupCommerce GroupType = "Commerce"
)

// student struct using GroupType
type Student struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	City      string    `json:"city"`
	Email     string    `json:"email"`
	Group     GroupType `json:"group"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	TutionFee float64   `json:"tution_fee"`
	Enrolled  bool      `json:"enrolled"`
	Mentor    string    `json:"mentor"`
	Subjects  []string  `json:"subjects"`
}

// validate function
func (g GroupType) IsValid() bool {
	return g == GroupScience || g == GroupCommerce || g == GroupArts
}