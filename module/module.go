package module

import (
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"database/sql"
)

type Module struct {
	AuthService    service.AuthService
	ProfileService service.ProfileService
}

func AuthModule(db *sql.DB) *Module {
	authRepository := repository.NewAuthRepository(db)
	profileRepository := repository.NewProfileRepository(db)

	authService := service.NewAuthService(authRepository)
	profileService := service.NewProfileService(profileRepository)

	return &Module{
		AuthService:    authService,
		ProfileService: profileService,
	}
}
