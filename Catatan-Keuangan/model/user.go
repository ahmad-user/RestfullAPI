package model

type Users struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	CreateAt  string `json:"createat"`
	UpdatedAt string `json:"updateat"`
}
