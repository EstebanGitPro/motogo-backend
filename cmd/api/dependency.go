package api

import (
	"log"

	"github.com/EstebanGitPro/motogo-backend/config"
	"github.com/EstebanGitPro/motogo-backend/internal/domain/person"
	"github.com/EstebanGitPro/motogo-backend/internal/platform/notification"

	//"github.com/EstebanGitPro/motogo-backend/internal/platform/jwt"  // <- Importar JWT
	repo "github.com/EstebanGitPro/motogo-backend/internal/platform/person"
)

type Dependencies struct {
	PersonService person.Service
	config        *config.Config
}

func initDependencies() *Dependencies {

	cfg := config.MustLoadConfig()

	db, err := repo.GetDB(cfg.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	personRepo := repo.NewRepository(db)

	resendNotifier, err := notification.NewResendNotifier(cfg.Resend.APIKey, cfg.Resend.FromEmail)
	if err != nil {
		log.Fatalf("Error creating Resend notifier: %v", err)
	}

	//jwtGenerator := jwt.New(cfg.JWT.SecretKey)

	personService := person.NewService(personRepo, resendNotifier, nil, cfg)

	return &Dependencies{
		PersonService: personService,
		config:        cfg,
	}
}
