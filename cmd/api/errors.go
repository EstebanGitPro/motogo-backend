package api

import (
	"errors"
	"net/http"

	domain "github.com/EstebanGitPro/motogo-backend/internal/domain/person"
	"github.com/gin-gonic/gin"
)

var (
	ErrUnmarshalBody      = errors.New("error unmarshal request body")
	ErrValidationUser     = errors.New("error validation user")
	ErrInvalidJSONFormat  = errors.New("invalid JSON format")
	ErrSchemaValidation   = errors.New("schema validation failed")
	ErrInvalidToken       = errors.New("invalid or expired token")

)

type WebError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (h handler) HandleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrUnmarshalBody):
		c.JSON(http.StatusBadRequest, WebError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	case errors.Is(err, ErrValidationUser):
		c.JSON(http.StatusUnauthorized, WebError{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	case errors.Is(err, ErrInvalidJSONFormat):
		c.JSON(http.StatusBadRequest, WebError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	case errors.Is(err, ErrSchemaValidation):
		c.JSON(http.StatusBadRequest, WebError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	case errors.Is(err, ErrInvalidToken):
		c.JSON(http.StatusBadRequest, WebError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	case errors.Is(err, domain.ErrUserCannotSave):
		c.JSON(http.StatusFailedDependency, WebError{
			Status:  http.StatusFailedDependency,
			Message: err.Error(),
		})
		return
	case errors.Is(err, domain.ErrGettingUserByEmail):
		c.JSON(http.StatusFailedDependency, WebError{
			Status:  http.StatusFailedDependency,
			Message: err.Error(),
		})
		return
	case errors.Is(err, domain.ErrNotFoundUserByEmail):
		c.JSON(http.StatusNotFound, WebError{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	case errors.Is(err, domain.ErrUserCannotFound):
		c.JSON(http.StatusNotFound, WebError{
			Status:  http.StatusNotFound,
			Message: err.Error(),
		})
		return
	case errors.Is(err, domain.ErrUserCannotGet):
		c.JSON(http.StatusFailedDependency, WebError{
			Status:  http.StatusFailedDependency,
			Message: err.Error(),
		})
		return
	case errors.Is(err, domain.ErrDuplicateUser):
		c.JSON(http.StatusConflict, WebError{
			Status:  http.StatusConflict,
			Message: err.Error(),
		})
	case errors.Is(err, domain.ErrorEmailNotVerified):
		c.JSON(http.StatusConflict, WebError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	default:
		c.JSON(http.StatusInternalServerError, WebError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
}
