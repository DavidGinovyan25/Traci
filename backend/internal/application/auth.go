package application

type Register struct {
	Username   string
	FirstName  string
	SecondName string
	Email      string
	Password   string
}

type Login struct {
	Email    string
	Password string
}
