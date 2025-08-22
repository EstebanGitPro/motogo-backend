package api

import (
	"fmt"
	"html/template"

	"net/http"

	domain "github.com/EstebanGitPro/motogo-backend/internal/domain/person"
	"github.com/gin-gonic/gin"
)

type handler struct {
	service domain.Service
}

func New(service domain.Service) *handler {
	return &handler{
		service: service,
	}
}

func (h handler) renderErrorPage(c *gin.Context, title, content string) {
	data := ResponseEmail{
		Title:   title,
		Content: template.HTML(fmt.Sprintf("<p>%s</p>", content)),
	}
	c.HTML(http.StatusBadRequest, "response.html", data)
}

func (h handler) handlerEmailVerificationError(c *gin.Context, err error) {
	switch err {
	case domain.ErrTokenExpired:
		h.renderErrorPage(c, "Token Expirado", "El enlace de verificación ha expirado. Por favor, solicita uno nuevo.")
	case domain.ErrTokenAlreadyUsed:
		h.renderErrorPage(c, "Token Ya Utilizado", "Este enlace ya ha sido utilizado anteriormente.")
	case domain.ErrTokenNotFound:
		h.renderErrorPage(c, "Token no válido", "El enlace de verificación no es válido.")
	default:
		h.HandleError(c, ErrValidationUser)
	}
}

func (h handler) GetPersonByEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Param("email")

		person, err := h.service.GetPersonByEmail(email)
		if err != nil {
			h.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, person)
	}
}

func (h handler) GetByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")

		person, err := h.service.GetByID(id)
		if err != nil {
			h.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, person)
	}
}

func (h handler) CheckEmailStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.Query("email")
		if email == "" {
			h.HandleError(c, domain.ErrorEmailNotVerified)
			return
		}

		person, err := h.service.GetPersonByEmail(email)
		if err != nil {
			h.HandleError(c, err)
			return
		}

		response := PersonEmailVerifiedResponse{
			EmailVerified: person.EmailVerified,
		}
		c.JSON(http.StatusOK, response)
	}
}

func (h handler) Save() func(c *gin.Context) {
	return func(c *gin.Context) {

		var personRequest PersonRequest
		if err := c.ShouldBindJSON(&personRequest); err != nil {
			h.HandleError(c, ErrInvalidJSONFormat)
			return
		}

		person, err := h.service.Save(personRequest.ToDomain())
		if err != nil {

			switch err {
			case domain.ErrDuplicateUser:
				h.HandleError(c, domain.ErrDuplicateUser)
			case domain.ErrUserCannotSave:
				h.HandleError(c, domain.ErrUserCannotSave)
			default:
				h.HandleError(c, domain.ErrUserCannotSave)
			}
			return
		}

		response := PersonResponse{
			ID:                  person.ID,
			IdentityNumber:      person.IdentityNumber,
			FirstName:           person.FirstName,
			LastName:            person.LastName,
			SecondLastName:      *person.SecondLastName,
			Email:               person.Email,
			PhoneNumber:         person.PhoneNumber,
			EmailVerified:       person.EmailVerified,
			PhoneNumberVerified: person.PhoneNumberVerified,
			Role:                person.Role,
		}

		c.JSON(http.StatusCreated, response)
	}
}

func (h handler) VerifyEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Param("token")
		if token == "" || len(token) < 20 {
			h.HandleError(c, ErrInvalidToken)
			return
		}

		err := h.service.VerifyEmailByToken(token)
		if err != nil {
			h.handlerEmailVerificationError(c, err)
			return
		}

		data := ResponseEmail{
			Title:   "Bienvenido a MotoGo",
			Content: template.HTML(`<p>Tu correo ha sido verificado exitosamente. Ahora puedes disfrutar de todas las funcionalidades.</p>`),
		}
		c.HTML(http.StatusOK, "response.html", data)
	}
}

func (h handler) Login() func(c *gin.Context) {
	return func(c *gin.Context) {
		var personLogin PersonLogin
		if err := c.ShouldBindJSON(&personLogin); err != nil {
			h.HandleError(c, ErrInvalidJSONFormat)
			return
		}

		person, token, err := h.service.Login(personLogin.ToDomain())
		if err != nil {
			if err == domain.ErrorEmailNotVerified {
				h.HandleError(c, domain.ErrorEmailNotVerified)
				return
			}
			h.HandleError(c, ErrValidationUser)
			return
		}

		response := LoginResponse{
			ID:                  person.ID,
			IdentityNumber:      person.IdentityNumber,
			FirstName:           person.FirstName,
			LastName:            person.LastName,
			SecondLastName:      *person.SecondLastName,
			Email:               person.Email,
			PhoneNumber:         person.PhoneNumber,
			EmailVerified:       person.EmailVerified,
			PhoneNumberVerified: person.PhoneNumberVerified,
			Role:                person.Role,
			Token:               token,
		}

		c.JSON(http.StatusOK, response)
	}

}

func (h *handler) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		var req PersonUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			h.HandleError(c, ErrInvalidJSONFormat)
			return
		}

		if err := h.service.Update(id, req.ToDomain()); err != nil {
			h.HandleError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func (h handler) SendPasswordRecoveryEmail() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request SendPasswordRecoveryRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			h.HandleError(c, ErrInvalidJSONFormat)
			return
		}

		err := h.service.SendPasswordRecoveryEmail(request.Email)
		if err != nil {
			switch err {
			case domain.ErrUserCannotFound:
				h.HandleError(c, domain.ErrUserCannotFound)
			default:
				h.HandleError(c, err)
			}
			return
		}

		response := GenericResponse{
			Status:  "success",
			Message: "password recovery email sent successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (h handler) RecoveryPassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request ResetPasswordRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			h.HandleError(c, ErrInvalidJSONFormat)
			return
		}

		
		userID, err := h.service.GetUserIDFromRecoveryCode(request.Code)
		if err != nil {
			switch err {
			case domain.ErrVerificationTokenNotFound:
				h.HandleError(c, domain.ErrVerificationTokenNotFound)
			case domain.ErrTokenExpired:
				h.HandleError(c, domain.ErrTokenExpired)
			case domain.ErrTokenAlreadyUsed:
				h.HandleError(c, domain.ErrTokenAlreadyUsed)
			default:
				h.HandleError(c, err)
			}
			return
		}

		err = h.service.VerifyPasswordRecoveryByCode(request.Code)
		if err != nil {
			switch err {
			case domain.ErrVerificationTokenNotFound:
				h.HandleError(c, domain.ErrVerificationTokenNotFound)
			case domain.ErrTokenExpired:
				h.HandleError(c, domain.ErrTokenExpired)
			case domain.ErrTokenAlreadyUsed:
				h.HandleError(c, domain.ErrTokenAlreadyUsed)
			default:
				h.HandleError(c, err)
			}
			return
		}

		err = h.service.RecoveryPassword(userID, request.NewPassword)
		if err != nil {
			h.HandleError(c, err)
			return
		}

		response := GenericResponse{
			Status:  "success",
			Message: "Password recovery completed successfully",
		}

		c.JSON(http.StatusOK, response)
	}
}
