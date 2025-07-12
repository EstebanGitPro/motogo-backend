package schema

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/EstebanGitPro/motogo-backend/tools/utils"
	"github.com/kaptinlin/jsonschema"
)

type Validators struct {
	FileReader FileReaderInterface
	LoginValidator *jsonschema.Schema
	RegisterValidator *jsonschema.Schema
}	

type FileReaderInterface interface {
	ReadJsonSchema(resourcePath string) ([]byte, error)
}

type DefaultFileReader struct{}

func (f *DefaultFileReader) ReadJsonSchema(resourcePath string) ([]byte, error) {
	
	root, err := utils.FindModuleRoot()

	if err != nil {
		return nil, err
	}

	data, err := os.Open(filepath.Join(root, "internal/platform/schema/json_schemas", resourcePath))
	if err != nil {
		return nil, err
	}
	defer data.Close()

	return io.ReadAll(data)
	
}

func NewValidator(fileReader FileReaderInterface) (*Validators, error) {
	validator := &Validators{
		FileReader: fileReader,
	}

	register, err := validator.createSchema("register_schema.json")
	if err != nil {
		return nil, err
	}

	validator.RegisterValidator = register

	return validator, nil

}

func (v *Validators) createSchema(resourcePath string) (*jsonschema.Schema, error) {
	compiler := jsonschema.NewCompiler()
	compiler.AssertFormat = true
	schemaJSON, err := v.FileReader.ReadJsonSchema(resourcePath) 
	if err != nil {
		return nil, errors.New("error reading json schema" + err.Error())
	}

	if schemaJSON == nil {
		return nil, errors.New("SchemaJSON is nil or empty")	
	}

	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		return nil, errors.New("error compiling schema" + err.Error())
	}

	return schema, nil
}


