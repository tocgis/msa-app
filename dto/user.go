package dto

type LoginDto struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Type     int    `json:"type"`
}

type AuthDto struct {
	Token string `json:"token"`
}

type UserDto struct {
	Id    int64  `json:"id"`
	Name  string `json:"Name"`
	Phone string `json:"Phone"`
}

type RegisterDto struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Code     string `json:"code"`
	Type     int    `json:"type"`
}
