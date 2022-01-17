package main

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	Name               string `json:"name, omitempty"`
	Lastname           string `json:"lastname, omitempty"`
	ExamID             int    `json:"exam, omitempty"`
	Accessed           bool   `json:"accessed, omitempty"`
	Submitted          bool   `json:"submitted, omitempty"`
	RegistrationNumber int    `json:"registration_number, omitempty"`
	jwt.StandardClaims
}

func NewCustomClaims(name, lastname string, examID int, expiration int64) CustomClaims {
	token := CustomClaims{
		Name:     name,
		Lastname: lastname,
		ExamID:   examID,
		Accessed: false,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration,
			Issuer:    "vano-jwt-teachers",
		},
	}
	return token
}

func NewSignedToken(claim CustomClaims) (string, error) {
	//unsigned token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	//sign the token
	return token.SignedString([]byte(conf.Secret))
}

func ParseToken(t string) (CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		t,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.Secret), nil
		},
	)
	if err != nil {
		return CustomClaims{}, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return CustomClaims{}, errors.New("can't parse claims")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return CustomClaims{}, errors.New("jwt is expired")
	}
	return *claims, nil
}
