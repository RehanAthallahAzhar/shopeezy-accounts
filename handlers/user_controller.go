package handlers

import (
	"net/http"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/pkg/errors"
	"github.com/labstack/echo/v4"
)

func (a *API) RegisterUser(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.UserRegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.ErrInvalidRequestPayload)
	}

	err := a.UserService.Register(ctx, &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondSuccess(c, http.StatusCreated, MsgUserCreated, nil)
}

func (a *API) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var req models.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, errors.ErrInvalidRequestPayload)
	}

	res, err := a.UserService.Login(ctx, &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondSuccess(c, http.StatusOK, MsgLogin, res)
}

func (a *API) Logout(c echo.Context) error {
	ctx := c.Request().Context()

	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, errors.ErrTokenNotFound)
	}

	if err := a.UserService.Logout(ctx, authHeader); err != nil {
		return handleServiceError(c, err)
	}

	return respondSuccess(c, http.StatusOK, MsgLogout, nil)
}

func (a *API) GetAllUsers(c echo.Context) error {
	ctx := c.Request().Context()

	res, err := a.UserService.GetAllUsers(ctx)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondWithUsers(c, res, err)
}

func (a *API) GetUserById(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := extractUserID(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, errors.ErrInvalidUserSession)
	}

	res, err := a.UserService.GetUserById(ctx, id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondWithUser(c, res, MsgUserRetrieved, err)
}

func (a *API) GetUserProfile(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := extractUserID(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, errors.ErrInvalidUserSession)
	}

	res, err := a.UserService.GetUserById(ctx, id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondWithUser(c, res, MsgUserRetrieved, err)
}

func (api *API) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := extractUserID(c)
	if err != nil {
		return respondError(c, http.StatusUnauthorized, errors.ErrInvalidUserSession)
	}

	var req models.UserUpdateRequest
	if err := c.Bind(&req); err != nil {
		return respondError(c, http.StatusBadRequest, errors.ErrInvalidRequestPayload)
	}

	res, err := api.UserService.UpdateUser(ctx, id, &req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondWithUser(c, res, MsgUserUpdated, err)
}

func (api *API) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := helpers.GetIDFromPathParam(c, "id")
	if err != nil {
		return respondError(c, http.StatusBadRequest, err)
	}

	err = api.UserService.DeleteUser(ctx, id)
	if err != nil {
		return handleServiceError(c, err)
	}

	return respondSuccess(c, http.StatusOK, MsgUserDeleted, nil)
}
