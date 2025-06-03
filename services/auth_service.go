package services

import (
	"errors"

	"github.com/donajivt/go-auth-service/db"
	"github.com/donajivt/go-auth-service/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Interface
type AuthService interface {
	Register(req models.RegistrationRequestDto) (string, error)
	Login(req models.LoginRequestDto) (models.LoginResponseDto, error)
	AssignRole(email, roleName string) error
}

type authService struct {
	jwt JwtService
}

func NewAuthService(jwt JwtService) AuthService {
	return &authService{jwt: jwt}
}

// Helpers
func hashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hash, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}

// Registro de usuario
func (s *authService) Register(req models.RegistrationRequestDto) (string, error) {
	var existing models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return "", errors.New("email already registered")
	}

	hashedPwd, _ := hashPassword(req.Password)

	// Verificar o crear rol
	var role models.Role
	if err := db.DB.Where("name = ?", req.Role).First(&role).Error; err != nil {
		role = models.Role{Name: req.Role}
		if err := db.DB.Create(&role).Error; err != nil {
			return "", err
		}
	}

	newUser := models.User{
		ID:          uuid.New(),
		Email:       req.Email,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPwd,
		RoleID:      role.ID,
	}

	if err := db.DB.Create(&newUser).Error; err != nil {
		return "", err
	}

	return newUser.ID.String(), nil
}

// Login de usuario
func (s *authService) Login(req models.LoginRequestDto) (models.LoginResponseDto, error) {
	var user models.User
	if err := db.DB.Preload("Role").Where("email = ?", req.UserName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LoginResponseDto{}, errors.New("invalid credentials")
		}
		return models.LoginResponseDto{}, err
	}

	if !checkPassword(user.Password, req.Password) {
		return models.LoginResponseDto{}, errors.New("invalid credentials")
	}

	token, _ := s.jwt.GenerateToken(user, user.Role)

	return models.LoginResponseDto{
		User: models.UserDto{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
		},
		Token: token,
	}, nil
}

// Asignar rol
func (s *authService) AssignRole(email, roleName string) error {
	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	var role models.Role
	if err := db.DB.FirstOrCreate(&role, models.Role{Name: roleName}).Error; err != nil {
		return err
	}

	user.RoleID = role.ID
	return db.DB.Save(&user).Error
}
