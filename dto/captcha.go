package dto


type SmsDto struct {
	Type string `json:"type"`
	Phone string `json:"phone"`
	Code string `json:"code"`
}

type CaptchaDto struct {
	X int `json:"x"`
	Y int `json:"y"`
	P int  `json:"p"`

}