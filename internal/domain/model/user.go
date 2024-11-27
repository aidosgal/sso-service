package model

type User struct {
	Id       int64
	Phone    string
	Name     string
	PassHash []byte
}
