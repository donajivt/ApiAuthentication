package services

import (
	"time"

	"github.com/donajivt/go-auth-service/config"
	"github.com/donajivt/go-auth-service/models"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateToken(user models.User, roles []string) (string, error)
}

type jwtService struct{}

func NewJwtService() JwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(user models.User, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"roles": roles,
		"iss":   config.Cfg.JwtOptions.Issuer,
		"aud":   config.Cfg.JwtOptions.Audience,
		"exp":   time.Now().AddDate(0, 0, 7).Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString([]byte(config.Cfg.JwtOptions.Secret))
}
