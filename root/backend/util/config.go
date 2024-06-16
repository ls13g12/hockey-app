package util

import "crypto/ed25519"

type Config struct {
	Addr 								string
	TokenPublicKey 			ed25519.PublicKey
	TokenPrivateKey 		ed25519.PrivateKey
	Dsn  								string
	Mode 								string
}
