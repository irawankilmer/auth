package service

import (
	"auth_service/internal/dto/request"
	"auth_service/internal/model"
	"auth_service/internal/repository"
	"context"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/gogaruda/pkg/utils"
)

type AuthService interface {
	Register(ctx context.Context, req request.RegisterRequest) error
	Login(req request.LoginRequest) (string, error)
	Logout(userID string) error
}

type authService struct {
	repo repository.AuthRepository
}

func NewAuthService(a repository.AuthRepository) AuthService {
	return &authService{repo: a}
}

func (s *authService) Register(ctx context.Context, req request.RegisterRequest) error {
	roleSet := make(map[string]bool)
	for _, r := range req.Roles {
		if roleSet[r] {
			msg := fmt.Sprintf("duplikat role ditemukan: %s", r)
			return apperror.New(apperror.CodeBadRequest, msg, nil)
		}
		roleSet[r] = true
	}

	usernameExists, err := s.repo.IsUsernameExists(req.Username)
	if err != nil {
		return err
	}
	if usernameExists {
		return apperror.New(apperror.CodeUsernameConflict, "username sudah terdaftar", err)
	}

	emailExists, err := s.repo.IsEmailExists(req.Email)
	if err != nil {
		return err
	}
	if emailExists {
		return apperror.New(apperror.CodeEmailConflict, "email sudah terdaftar", err)
	}

	validRoles, err := s.repo.CheckRoles(req.Roles)
	if err != nil {
		return err
	}

	passHash, err := utils.GenerateHash(req.Password)
	if err != nil {
		return apperror.New(apperror.CodeInternalError, "gagal generate password", err)
	}

	user := &model.UserModel{
		ID:       utils.NewULID(),
		Username: req.Username,
		Email:    req.Email,
		Password: passHash,
		Roles:    validRoles,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(req request.LoginRequest) (string, error) {
	user, err := s.repo.IdentifierCheck(req.Identifier)
	if err != nil || !utils.CompareHash(user.Password, req.Password) {
		return "", err
	}

	var roleNames []string
	for _, r := range user.Roles {
		roleNames = append(roleNames, r.Name)
	}

	newVersion, err := s.repo.UpdateTokenVersion(user.ID)
	if err != nil {
		return "", err
	}

	token, err := utils.GenerateJWT(user.ID, newVersion, roleNames)
	if err != nil {
		return "", apperror.New(apperror.CodeInternalError, "gagal membuat jwt", err)
	}

	return token, nil
}

func (s *authService) Logout(userID string) error {
	_, err := s.repo.UpdateTokenVersion(userID)
	if err != nil {
		return err
	}

	return nil
}
