package jwt_handler

// DecodeUserToken 用户token信息
type DecodeUserToken struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Iat  int64  `json:"iat"`
	Iss  string `json:"iss"`
	Exp  int64  `json:"exp"`
}

// DecodeUserRefreshToken 用户refreshToken信息
type DecodeUserRefreshToken struct {
	ID  uint   `json:"id"`
	Iat int64  `json:"iat"`
	Iss string `json:"iss"`
	Exp int64  `json:"exp"`
}
