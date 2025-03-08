package entity

//MIGHTDO : setting max and stuff in the models?

type User struct {
	ID       int64
	Username string
	Password string
	Balance  int64
}
