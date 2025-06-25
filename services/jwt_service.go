package services

import (
	"time"

	"github.com/donajivt/go-auth-service/config"
	"github.com/donajivt/go-auth-service/models"
	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateToken(user models.User, role models.Role) (string, error)
}

type jwtCustomService struct{}

func NewJwtService() JwtService {
	return &jwtCustomService{}
}

func (s *jwtCustomService) GenerateToken(user models.User, role models.Role) (string, error) {
	claims := jwt.MapClaims{
		"sub":         user.ID.String(),
		"name":        user.Name,
		"email":       user.Email,
		"phoneNumber": user.PhoneNumber,
		"role":        role.Name,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
		"iss":         config.Cfg.JwtOptions.Issuer,
		"aud":         config.Cfg.JwtOptions.Audience,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Cfg.JwtOptions.Secret))
}
