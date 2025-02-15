package types

type UserCreate struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	Birthdate string `json:"birthdate"`
}

type UserLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDAO struct {
	ID       int64  `db:"user_id"`
	Username string `db:"username"`
	Name     string `db:"name"`
	Password string `db:"password"`
}
