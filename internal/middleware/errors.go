package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)


var (
	ErrInvalidJSONFormat       = errors.New("invalid JSON format")
	ErrUnmarshalBody           = errors.New("failed to process request body")
	ErrSchemaValidation        = errors.New("validation failed")
	ErrInternalServer          = errors.New("internal server error")
	ErrModuleRootNotFound      = errors.New("could not find module root")
	ErrSchemaFileNotFound      = errors.New("schema file not found")
	ErrSchemaFileRead          = errors.New("failed to read schema file")
	ErrSchemaCompilation       = errors.New("failed to compile JSON schema")
	ErrSchemaEmpty             = errors.New("JSON schema is empty or null")
	ErrValidatorInitFailed     = errors.New("validator initialization failed")
	ErrValidationUserFailed    = errors.New("user validation failed")
	ErrValidationUserNotFound  = errors.New("user not found")
	ErrValidationUserAlreadyExists = errors.New("user already exists")
)


func ValidateError(c *gin.Context, err error, details interface{}, statusCode int) {
	if details == nil {
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	
	detailsMap, ok := details.(map[string]interface{})
	if !ok {
		c.JSON(statusCode, gin.H{
			"error":   err.Error(),
			"details": details,
		})
		c.Abort()
		return
	}
	
	fieldErrors := make(map[string]string)
	

	if detailsList, exists := detailsMap["details"].([]interface{}); exists {
		for _, item := range detailsList {
			fieldDetail, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			
			if valid, exists := fieldDetail["valid"].(bool); exists && !valid {
				
				path := ""
				if instanceLoc, exists := fieldDetail["instanceLocation"].(string); exists {
					path = instanceLoc
					if len(path) > 0 && path[0] == '/' {
						path = path[1:]
					}
				}
				
				if errorsMap, exists := fieldDetail["errors"].(map[string]interface{}); exists && path != "" {
					for _, msg := range errorsMap {
						if strMsg, ok := msg.(string); ok {
							fieldErrors[path] = strMsg
						}
					}
				}
			}
		}
	}
	
	c.JSON(statusCode, gin.H{
		"error":   err.Error(),
		"invalid": fieldErrors,
	})
	c.Abort()
}





