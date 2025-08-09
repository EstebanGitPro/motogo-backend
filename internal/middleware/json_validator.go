package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	schema "github.com/EstebanGitPro/motogo-backend/internal/platform/schema"
	"github.com/gin-gonic/gin"
	"github.com/kaptinlin/jsonschema"
)

type Builder struct {
	Validators *schema.Validators
    isLogin    bool 
}


func NewMiddlewareValidator(validators *schema.Validators) *Builder {
	
	return &Builder{
		Validators: validators,
	}
}

func (b *Builder) WithValidateLogin() gin.HandlerFunc {
    b.isLogin = true
	return b.jsonValidator(b.Validators.LoginValidator)
}


func (b *Builder) WithValidateRegister() gin.HandlerFunc {
    b.isLogin = false
	return b.jsonValidator(b.Validators.RegisterValidator )
}

func (b *Builder) WithValidateUpdate() gin.HandlerFunc {
    b.isLogin = false
	return b.jsonValidator(b.Validators.UpdateValidator) 
}

func (b *Builder) jsonValidator(schema *jsonschema.Schema) gin.HandlerFunc {
    return func(c *gin.Context) {

        bodyBytes, err := io.ReadAll(c.Request.Body)
        if err != nil { 
            ValidateError(c, ErrUnmarshalBody, nil, http.StatusBadRequest)
            return
        }

      
        c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

      
        var data map[string]interface{}
        if err := json.Unmarshal(bodyBytes, &data); err != nil {
            ValidateError(c, ErrInvalidJSONFormat, nil, http.StatusBadRequest)
            return
        }

    
        result := schema.Validate(data)
        if !result.IsValid() {
            details, err := json.MarshalIndent(result.ToList(), "", "  ")
            if err != nil {
                ValidateError(c, ErrUnmarshalBody, nil, http.StatusBadRequest)
                return
            }
            

            var formattedErrors interface{}
            if err := json.Unmarshal(details, &formattedErrors); err != nil {
                ValidateError(c, ErrUnmarshalBody, nil, http.StatusBadRequest)
                return
            }
            
            ValidateError(c, ErrSchemaValidation, formattedErrors, http.StatusBadRequest)
            return
        }

        
        c.Next()
    }
}
