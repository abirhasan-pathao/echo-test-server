package usecase

import (
	"echo-server/app/utils"
)

type CookieUsecase struct{}

func NewCookieUsecase() *CookieUsecase {

	return &CookieUsecase{}
}

func (u *CookieUsecase) GetCookie() (string, error) {
	token, err := utils.GenerateJWT()
	if err != nil {
		return "", err
	}

	return token, err
}
