package main

import (
	"github.com/pbulteel/homebox-justfind/backend/internal/core/services"
	"github.com/pbulteel/homebox-justfind/backend/internal/core/services/reporting/eventbus"
	"github.com/pbulteel/homebox-justfind/backend/internal/data/ent"
	"github.com/pbulteel/homebox-justfind/backend/internal/data/repo"
	"github.com/pbulteel/homebox-justfind/backend/internal/sys/config"
	"github.com/pbulteel/homebox-justfind/backend/pkgs/mailer"
)

type app struct {
	conf        *config.Config
	mailer      mailer.Mailer
	db          *ent.Client
	repos       *repo.AllRepos
	services    *services.AllServices
	bus         *eventbus.EventBus
	authLimiter *authRateLimiter
}

func new(conf *config.Config) *app {
	s := &app{
		conf: conf,
	}

	s.mailer = mailer.Mailer{
		Host:     s.conf.Mailer.Host,
		Port:     s.conf.Mailer.Port,
		Username: s.conf.Mailer.Username,
		Password: s.conf.Mailer.Password,
		From:     s.conf.Mailer.From,
	}

	s.authLimiter = newAuthRateLimiter(s.conf.Auth.RateLimit)

	return s
}
