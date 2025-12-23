package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"be-job-portal/internal/config"
	"be-job-portal/internal/domain"
	"be-job-portal/pkg/utils"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type authUsecase struct {
	userRepo domain.UserRepository
	cfg      config.Config
	oauthCfg *oauth2.Config
}

func NewAuthUsecase(userRepo domain.UserRepository, cfg config.Config) domain.AuthUsecase {
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.GoogleRedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	return &authUsecase{
		userRepo: userRepo,
		cfg:      cfg,
		oauthCfg: oauthCfg,
	}
}

func (u *authUsecase) Register(ctx context.Context, email, password, role string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		Provider: "local",
	}

	return u.userRepo.Create(ctx, user)
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if user.Provider != "local" {
		return "", fmt.Errorf("please login with %s", user.Provider)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateToken(user.ID, user.Role, u.cfg.JWTSecret)
}

func (u *authUsecase) GoogleLogin(ctx context.Context, code string) (string, error) {
	token, err := u.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return "", err
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", err
	}

	user, err := u.userRepo.GetByEmail(ctx, googleUser.Email)
	if err != nil {
		user = &domain.User{
			Email:    googleUser.Email,
			Role:     "SEEKER",
			Provider: "google",
		}
		if err := u.userRepo.Create(ctx, user); err != nil {
			return "", err
		}
	}

	return utils.GenerateToken(user.ID, user.Role, u.cfg.JWTSecret)
}
