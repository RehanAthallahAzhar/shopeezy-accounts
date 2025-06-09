package handlers

import (
	"errors"
	"net/http"

	"github.com/RehanAthallahAzhar/shopeezy-accounts/helpers"
	"github.com/RehanAthallahAzhar/shopeezy-accounts/models"
	"github.com/labstack/echo/v4"
)

func (api *API) FindAllUsers(c echo.Context) error {
	ctx := c.Request().Context()

	// Inisialisasi slice untuk menampung data user
	var users []models.User

	users, err := api.UserRepo.ReadAllProducts(ctx)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	if len(users) == 0 {
		return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User list not found!"})
	}

	// Kirimkan response sukses dengan data user
	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Lists Data Users",
		Data:    users,
	})
}

func (api *API) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	//struct user request
	var req = models.UserCreateRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
	}

	// Inisialisasi user baru
	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	err := api.UserRepo.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	// Kirimkan response sukses
	return c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "User created successfully",
		Data: models.UserResponse{
			Id:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})

}

func (api *API) FindUserById(c echo.Context) error {
	ctx := c.Request().Context()

	// Ambil ID user dari parameter URL
	id := c.Param("id")

	// Inisialisasi user
	var user models.User
	user, err := api.UserRepo.FindUserById(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	// Kirimkan response sukses dengan data user
	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User Founded",
		Data: models.UserResponse{
			Id:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})

}

func (api *API) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	userId := c.Get("userID").(string)

	// Inisialisasi user
	var user models.User
	user, err := api.UserRepo.FindUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	//struct user request
	var req = models.UserUpdateRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid JSON format"})
	}

	user = models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: helpers.HashPassword(req.Password),
	}

	err = api.UserRepo.UpdateUser(ctx, userId, &user)
	if err != nil {

		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	// Kirimkan response sukses
	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User updated successfully",
		Data: models.UserResponse{
			Id:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func (api *API) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()

	// Ambil ID user dari parameter URL
	id := c.Param("id")

	_, err := api.UserRepo.FindUserById(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	err = api.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found!"})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to retrieve user"})
	}

	// Kirimkan response sukses
	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "User deleted successfully",
	})
}
