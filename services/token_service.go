package services

import (
	"chat/config"
	"chat/src/api/pkg/logging"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	Logger logging.Logger
	Cfg    *config.Config
}

type TokenDto struct {
	UserID       int
	FirstName    string
	LastName     string
	UserName     string
	Email        string
	MobileNumber string
	Roles        []string
}

type TokenDetail struct {
	AccessToken                string `json:"accessToken"`
	RefreshToken               string `json:"refreshToken"`
	AccessTokenExpirationTime  int64  `json:"accessTokenExpirationTime"`
	RefreshTokenExpirationTime int64  `json:"refreshTokenExpirationTime"`
}

func NewTokenService(cfg *config.Config) *TokenService {
	logger := logging.NewLogger(cfg)
	u := &TokenService{
		Logger: logger,
		Cfg:    cfg,
	}
	return u
}

func (tokenService *TokenService) GenerateToken(tokenDto *TokenDto) (*TokenDetail, error) {
	token := &TokenDetail{}
	token.AccessTokenExpirationTime = time.Now().Add(tokenService.Cfg.JWT.AccessTokenExpireDuration * time.Minute).Unix()
	token.RefreshTokenExpirationTime = time.Now().Add(tokenService.Cfg.JWT.RefreshTokenExpireDuration * time.Minute).Unix()

	dto := jwt.MapClaims{} //access token claims

	dto["id"] = tokenDto.UserID
	dto["firstName"] = tokenDto.FirstName
	dto["lastName"] = tokenDto.LastName
	dto["userName"] = tokenDto.UserName
	dto["email"] = tokenDto.Email
	dto["mobileNumber"] = tokenDto.MobileNumber
	dto["roles"] = tokenDto.Roles
	dto["exp"] = token.AccessTokenExpirationTime

	atc := jwt.NewWithClaims(jwt.SigningMethodHS256, dto)
	var err error
	token.AccessToken, err = atc.SignedString([]byte(tokenService.Cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}

	dtor := jwt.MapClaims{} //refresh token claims
	dtor["id"] = tokenDto.UserID
	dtor["exp"] = token.RefreshTokenExpirationTime

	rtc := jwt.NewWithClaims(jwt.SigningMethodHS256, dtor)

	token.RefreshToken, err = rtc.SignedString([]byte(tokenService.Cfg.JWT.Secret))

	if err != nil {
		return nil, err
	}
	return token, nil
}

func (tokenService *TokenService) ValidateToken(token string) (*jwt.Token, error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenService.Cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}
	return tokenObj, nil
}

func (tokenService *TokenService) GetClaims(token string) (mapClaims map[string]interface{}, err error) {
	mapClaims = map[string]interface{}{}
	str := "claims"
	jwtToken, err := tokenService.ValidateToken(token)

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("cant get %v", str)
	}
	for i, v := range claims {
		mapClaims[i] = v
	}
	return mapClaims, nil
}
