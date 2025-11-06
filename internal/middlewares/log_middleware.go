package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware adalah middleware untuk mencatat setiap request yang masuk.
func LoggingMiddleware(log *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Logika ini akan berjalan SEBELUM handler dieksekusi
			req := c.Request()
			res := c.Response()

			log.WithFields(logrus.Fields{
				"request_id": res.Header().Get(echo.HeaderXRequestID),
				"ip_address": c.RealIP(),
				"method":     req.Method,
				"path":       req.URL.Path,
			}).Info("Received new request")

			// Teruskan request ke handler berikutnya
			err := next(c)

			// Logika ini bisa berjalan SETELAH handler dieksekusi (opsional)
			// Misalnya, untuk mencatat status code
			log.WithFields(logrus.Fields{
				"request_id": res.Header().Get(echo.HeaderXRequestID),
				"status":     res.Status,
			}).Info("Finished processing request")

			return err
		}
	}
}
