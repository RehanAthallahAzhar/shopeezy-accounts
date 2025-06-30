package helpers

import (
	"github.com/RehanAthallahAzhar/shopeezy-accounts/pkg/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func GenerateNewUUID() string {
	return uuid.New().String()
}

func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func GetIDFromPathParam(c echo.Context, key string) (string, error) {
	val := c.Param(key)
	if val == "" || !isValidUUID(val) {
		return "", errors.ErrInvalidRequestPayload
	}
	return val, nil
}

func GetFromPathParam(c echo.Context, key string) (string, error) {
	val := c.Param(key)
	if val == "" {
		return "", errors.ErrInvalidRequestPayload
	}
	return val, nil
}
