package types

type UserCreate struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	Birthdate string `json:"birthdate"`
}

type UserLoginDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDAO struct {
	ID       int64  `db:"user_id"`
	Email    string `db:"email"`
	Name     string `db:"name"`
	Password string `db:"password"`
}
