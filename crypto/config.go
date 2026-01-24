package crypto

type Cfg struct {
	JWT_EXP int    // minute
	JWT_KEY string // secret
}

var Config *Cfg
