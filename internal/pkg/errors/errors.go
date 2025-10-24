package errors

import "errors"

var (
	// validation
	ErrInvalidRequestPayload = errors.New("invalid request payload")
	ErrInvalidQuery          = errors.New("invalid query parameters")
	ErrInvalidUserSession    = errors.New("invalid user session")
	ErrInvalidToken          = errors.New("invalid token")
	ErrExpiredToken          = errors.New("expired token")
	ErrInvalidTokenRole      = errors.New("expired token role")
	ErrFailedToGenerateToken = errors.New("failed to generate token")
	ErrInvalidTokenFormat    = errors.New("invalid token format")
	ErrMissingJTI            = errors.New("missing JTI, cannot revoke token")
	ErrFailedToRevokeToken   = errors.New("failed to revoke token")
	ErrTokenNotFound         = errors.New("token not found")
	ErrInvalidCredentials    = errors.New("invalid credentials")

	// status
	ErrInternalServerError = errors.New("internal server error")

	// user
	ErrUserNotFound       = errors.New("user not found")
	ErrFailedToReadUser   = errors.New("failed to read user")
	ErrFailedToCreateUser = errors.New("failed to create user")
	ErrFailedToUpdateUser = errors.New("failed to update user")
	ErrFailedToDeleteUser = errors.New("failed to delete user")
)

/*

### 📌 Error Handling Best Practice per Layer

| Layer       	   | Gunakan `errors.New(...)`                            		| Gunakan `fmt.Errorf(...)`                            |
|-----------------|-------------------------------------------------------------|------------------------------------------------------|
| **Repo**        | Untuk error domain yang dikenal (contoh: `ErrUserNotFound`) | Untuk memberikan detail/debug error DB (ex: GORM)     |
| **Service**     | Untuk validasi input/domain (contoh: `ErrInvalidInput`)     | Untuk *wrap* error dari repository dengan konteks     |
| **Controller**  | Untuk membentuk respon HTTP user-friendly                | ❌ Tidak digunakan (cukup gunakan error dari service) |



*/
