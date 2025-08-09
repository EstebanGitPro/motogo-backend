package api

import (

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/internal/domain/person"
	"github.com/EstebanGitPro/motogo-backend/internal/domain/token"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/jwt"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/notification"
	repo "github.com/EstebanGitPro/motogo-backend/internal/platform/person"
)

type Dependencies struct {
	PersonService person.Service
	PersonRepo    person.Repository
	Notifier      person.Notifier
	Config        *config.Config
	JwtGenerator  token.Generator
}

func initDependencies() (*Dependencies, error) {


	cfg := config.MustLoadConfig()

	db, err := repo.GetDB(cfg.Database)
	if err != nil {
		return nil , err
	}

	personRepo := repo.NewRepository(db)

	resendNotifier, err := notification.NewResendNotifier(cfg.Resend.APIKey, cfg.Resend.FromEmail)
	if err != nil {
		return nil, err
	}

	jwtGenerator := jwt.New(cfg.JWT.SecretKey)

	personService := person.NewService(personRepo, resendNotifier, jwtGenerator, cfg)

	return &Dependencies{
		PersonService: personService,
		PersonRepo:    personRepo,
		Notifier:      resendNotifier,
		Config:        cfg,
		JwtGenerator:  jwtGenerator,
	}, nil
}
