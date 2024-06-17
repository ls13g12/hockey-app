package token

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	privateKey  ed25519.PrivateKey
	publicKey   ed25519.PublicKey
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (Maker, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key size: must be exactly %d bytes", ed25519.PrivateKeySize)
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size: must be exactly %d bytes", ed25519.PublicKeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		privateKey: 	privateKey,
		publicKey: 		publicKey,
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Sign(maker.privateKey, payload, nil)
	return token, payload, err
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Verify(token, maker.publicKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
