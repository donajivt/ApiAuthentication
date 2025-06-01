package models

import "github.com/google/uuid"

// --- Requests ---

type RegistrationRequestDto struct {
	Email       string `json:"email" binding:"required,email"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password" binding:"required,min=6"`
	Role        string `json:"role"`
}

type LoginRequestDto struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// --- Responses ---

type UserDto struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber"`
}

type LoginResponseDto struct {
	User  UserDto `json:"user"`
	Token string  `json:"token"`
}

type ResponseDto struct {
	Result    interface{} `json:"result,omitempty"`
	IsSuccess bool        `json:"isSuccess"`
	Message   string      `json:"message"`
}
