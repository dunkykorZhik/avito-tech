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

type jwtCustomClaim struct {
	userID   int64
	username string
	jwt.RegisteredClaims
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
	userID := user.ID
	if err == errs.ErrNoUser {
		userID, err = s.createUser(ctx, username, password)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	} else if user.Password != password {
		//hash the password
		return "", errs.ErrNoUser

	}

	claims := &jwtCustomClaim{
		userID:   userID,
		username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) createUser(ctx context.Context, username, password string) (int64, error) {
	//hash the password

	return s.userRepo.CreateUser(ctx, username, password)

}

func (s *UserService) ParseToken(tokenString string) (int64, string, error) {
	claims := &jwtCustomClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signing method")
		}
		return []byte(s.signKey), nil
	})
	if err != nil {
		return 0, "", err
	}
	if !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}
	return claims.userID, claims.username, nil
}
