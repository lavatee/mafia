package mafia

type User struct {
	Id       int    `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Name     string `json:"name" db:"name"`
	Password string `json:"password" db:"password_hash"`
	Coins    int    `json:"coins" db:"coins"`
}
