package service

import (
	"context"
	"fmt"
	"time"

	errs "github.com/dunkykorZhik/avito-tech/internal/errors"
	"github.com/dunkykorZhik/avito-tech/internal/repo"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	userRepo repo.User
	signKey  string
	ttl      time.Duration
}

func NewUserService(userRepo repo.User, key string, ttl time.Duration) *UserService {
	return &UserService{
		userRepo: userRepo,
		signKey:  key,
		ttl:      ttl,
	}
}

func (s *UserService) GenerateToken(ctx context.Context, username string, password string) (string, error) {
	//get user ->if do not exist, call create
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()
	user, err := s.userRepo.GetUserByName(ctx, username)
	if err == errs.ErrNoUser {
		err = s.createUser(ctx, username, password)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else if user.Password != password {
		//hash the password
		return "", errs.ErrNoUser

	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
		Subject:   username,
	}

	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) createUser(ctx context.Context, username, password string) error {
	//hash the password

	return s.userRepo.CreateUser(ctx, username, password)

}

func (s *UserService) ParseToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signing method")
		}
		return []byte(s.signKey), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return "", fmt.Errorf("cannot parse token")
	}
	return claims.Subject, nil
}
