package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lvjp/wtf-go/internal/app/config"
	"github.com/lvjp/wtf-go/pkg/auth/token"
)

var errAuthenticationFailed = errors.New("authentication failed")

type Service interface {
	CreateToken(ctx context.Context, subject string) (*token.Token[*SessionData], error)
	LookupToken(ctx context.Context, id string) (*token.Token[*SessionData], error)
	RevokeToken(ctx context.Context, id string) error
}

func NewService(cfg config.Auth) Service {
	// TODO: allow selecting the store backend via configuration once a SQL database provider is implemented.
	m, err := token.NewManager[*SessionData](token.UUIDGenerator, token.NewMemoryStore())
	if err != nil {
		panic(fmt.Sprintf("auth.NewService: unexpected error: %v", err))
	}
	return &service{manager: m, tokenTTL: cfg.TokenTTL}
}

type service struct {
	manager  token.Manager[*SessionData]
	tokenTTL time.Duration
}

func (s *service) CreateToken(ctx context.Context, subject string) (*token.Token[*SessionData], error) {
	// TODO: implement a proper identification layer instead of this hardcoded subject.
	if subject != "admin@localhost" {
		return nil, errAuthenticationFailed
	}

	return s.manager.Create(ctx, s.tokenTTL, &SessionData{Subject: subject})
}

func (s *service) LookupToken(ctx context.Context, id string) (*token.Token[*SessionData], error) {
	return s.manager.Lookup(ctx, id)
}

func (s *service) RevokeToken(ctx context.Context, id string) error {
	return s.manager.Revoke(ctx, id)
}
