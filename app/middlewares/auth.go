package middlewares

import "github.com/labstack/echo/v5"

func DummyAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c *echo.Context) error {
		user, pass, ok := c.Request().BasicAuth()
		if !ok || user != "admin" || pass != "password" {
			return c.JSON(401, map[string]string{"error": "Unauthorized"})
		}

		return next(c)
	}
}
