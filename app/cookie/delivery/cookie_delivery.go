package delivery

import (
	"echo-server/app/cookie/usecase"
	"echo-server/app/utils"
	"net/http"

	"github.com/labstack/echo/v5"
)

type CookieController struct {
	cookieUsecase *usecase.CookieUsecase
}

func NewCookieController(cookieUsecase *usecase.CookieUsecase) *CookieController {

	return &CookieController{cookieUsecase: cookieUsecase}
}

func (ctrl *CookieController) GetCookie(c *echo.Context) error {
	existingCookie, err := c.Cookie("cookie")
	if err == nil && existingCookie != nil {
		isValid, err := utils.ValidateJWT(existingCookie.Value)
		if isValid && err == nil {
			return c.JSON(http.StatusOK, map[string]string{"message": "You just got a cookie, go get some milk!"})
		}
	}

	cookieToken, err := ctrl.cookieUsecase.GetCookie()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to bake cookie"})
	}

	cookie := new(http.Cookie)
	cookie.Name = "cookie"
	cookie.Value = cookieToken
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "Here is your cookie! Go get some milk!"})
}

func (ctrl *CookieController) GetMilk(c *echo.Context) error {
	existingCookie, err := c.Cookie("cookie")
	if err != nil || existingCookie == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "You don't have the cookie, What will you do with milk?"})
	}
	isValid, err := utils.ValidateJWT(existingCookie.Value)
	if !isValid || err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Your cookie expired, go get new cookie"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Here is your milk! Enjoy!"})
}
