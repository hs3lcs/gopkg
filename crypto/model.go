package crypto

type JwtPack struct {
	ISS  string `json:"iss"`
	IAT  int64  `json:"iat"`
	EXP  int64  `json:"exp"`
	UID  uint32 `json:"uid"`
	ORG  uint16 `json:"org"`
	Type uint8  `json:"type"`
	Role uint8  `json:"role"`
}

const CHAR_HASH string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Cfg struct {
	JWT_EXP int    // second
	JWT_KEY string // secret
}

var Config *Cfg
