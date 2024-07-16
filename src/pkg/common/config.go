package common

import "crypto/ed25519"

type ServerConfig struct {
	Addr            string
	TokenPublicKey  ed25519.PublicKey
	TokenPrivateKey ed25519.PrivateKey
	Dsn             string
	Mode            string
}

type TuiAppConfig struct {
	Dsn  string
	Mode string
}
