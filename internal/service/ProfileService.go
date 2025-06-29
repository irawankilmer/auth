package service

import (
	"auth_service/internal/dto/request"
	"auth_service/internal/dto/response"
	"auth_service/internal/repository"
	"context"
)

type ProfileService interface {
	Setting(req request.ProfileRequest) error
	Me(ctx context.Context, userID string) (*response.UserResponse, error)
}

type profileService struct {
	profileRepo repository.ProfileRepository
}

func NewProfileService(p repository.ProfileRepository) ProfileService {
	return &profileService{profileRepo: p}
}

func (s *profileService) Setting(req request.ProfileRequest) error {
	exists, err := s.profileRepo.IsUserExists(req.UserID)
	if err != nil {
		return err
	}

	profile := request.ProfileRequest{
		UserID:   req.UserID,
		FullName: req.FullName,
		Address:  req.Address,
		Gender:   req.Gender,
	}

	if exists {
		if err := s.profileRepo.Update(profile); err != nil {
			return err
		}
	} else {
		if err := s.profileRepo.Create(profile); err != nil {
			return err
		}
	}

	return nil
}

func (s *profileService) Me(ctx context.Context, userID string) (*response.UserResponse, error) {
	return s.profileRepo.Me(ctx, userID)
}
