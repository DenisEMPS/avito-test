package types

type UserCreate struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	Birthdate string `json:"birthdate" binding:"required"`
}

type UserLoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDAO struct {
	ID       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
