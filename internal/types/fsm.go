package types

type FSM interface {
	CreateUser(id int64, state string) bool

	GetStateUser(id int64) (string, bool)
	SetStateUser(id int64, state string) bool

	SetValue(id int64, param string, data string) bool
	GetValue(id int64, param string) (string, bool)
}
