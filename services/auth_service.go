package services

import (
	"errors"
	"strings"

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
	AssignRole(email, role string) error
}

type authService struct {
	jwt JwtService
}

func NewAuthService(jwt JwtService) AuthService {
	return &authService{jwt: jwt}
}

// ---- Helpers ----
func hashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(hash, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}

// ---- Methods ----

func (s *authService) Register(req models.RegistrationRequestDto) (string, error) {
	var existing models.User

	// Buscar si el usuario ya existe por email
	if err := db.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return "", errors.New("email already registered")
	}

	// Hashear la contraseña
	pw, _ := hashPassword(req.Password)

	userID := uuid.New()

	newUser := models.User{
		ID:          userID,
		Email:       req.Email,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    pw,
	}

	// Guardar el usuario primero
	if err := db.DB.Create(&newUser).Error; err != nil {
		return "", err
	}

	// Buscar o crear el rol
	role := models.Role{}
	if err := db.DB.Where("name = ?", req.Role).First(&role).Error; err != nil {
		// Si no existe, lo creamos
		role = models.Role{
			ID:   role.ID,
			Name: req.Role,
		}
		if err := db.DB.Create(&role).Error; err != nil {
			return "", err
		}
	}

	// Insertar en tabla user_roles
	IdMayusculas := strings.ToUpper(string(userID.String()))
	sql := `INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)`
	if err := db.DB.Exec(sql, IdMayusculas, role.ID).Error; err != nil {
		return "", err
	}

	return newUser.ID.String(), nil
}

func (s *authService) Login(req models.LoginRequestDto) (models.LoginResponseDto, error) {
	var user models.User
	dberr := db.DB.Preload("Roles").First(&user, "email = ?", req.UserName).Error
	if errors.Is(dberr, gorm.ErrRecordNotFound) {
		return models.LoginResponseDto{}, errors.New("invalid credentials")
	}
	if !checkPassword(user.Password, req.Password) {
		return models.LoginResponseDto{}, errors.New("invalid credentials")
	}
	var roleNames []string
	for _, r := range user.Roles {
		roleNames = append(roleNames, r.Name)
	}
	token, _ := s.jwt.GenerateToken(user, roleNames)
	return models.LoginResponseDto{
		User: models.UserDto{
			ID:          user.ID,
			Email:       user.Email,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
		Token: token,
	}, nil
}

func (s *authService) AssignRole(email, roleName string) error {
	var user models.User
	if err := db.DB.First(&user, "email = ?", email).Error; err != nil {
		return err
	}
	role := models.Role{Name: roleName}
	db.DB.FirstOrCreate(&role, role)
	return db.DB.Model(&user).Association("Roles").Append(&role)
}
