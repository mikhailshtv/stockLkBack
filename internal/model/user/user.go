package user

type User struct {
	Id           int
	Login        string
	passwordHash string
	FirstName    string
	LastName     string
	Email        string
}

func (user *User) GetPasswordHash() string {
	return user.passwordHash
}

func (user *User) SetPasswordHash(passwordHash string) {
	user.passwordHash = passwordHash
}
