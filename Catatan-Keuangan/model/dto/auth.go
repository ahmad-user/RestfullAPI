package dto

type AuthResponDto struct {
	Token string `json:"token"`
}

type AuthReqDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
