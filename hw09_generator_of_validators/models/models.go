package models

//go:generate go-validate $GOFILE

type UserRole string

// NOTE: Several struct specs in one type declaration are allowed.
<<<<<<< HEAD
=======

>>>>>>> 4949677b261b4a96c4613acd5edd3da796c9443f
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\w+@\w+\.\w+$"` //nolint
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
	}

	App struct {
		Version string `validate:"len:5"`
	}
)

type Token struct {
	Header    []byte
	Payload   []byte
	Signature []byte
}

type Response struct {
	Code int    `validate:"in:200,404,500"`
	Body string `json:"omitempty"`
}
