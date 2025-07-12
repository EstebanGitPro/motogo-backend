package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/EstebanGitPro/motogo-backend/tools/utils"
)

type Config struct {
	Environment string    `json:"environment"`
	Database    Database  `json:"database"`
	Server      Server    `json:"server"`
	Resend      Resend    `json:"resend"`
	JWT         JWTConfig `json:"jwt"`
	Verification Verification `json:"verification"`
}

type Verification struct {
	BaseURL string `json:"base_url"`
}

type Database struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	URL      string `json:"url,omitempty"`
	SSL      string `json:"ssl,omitempty"`
}

type Server struct {
	Port string `json:"port"`
	Host string `json:"host"`
}

type Resend struct {
	APIKey    string `json:"api_key"`
	FromEmail string `json:"from_email"`
}

type JWTConfig struct {
	SecretKey string `json:"secret_key"`
}

func LoadConfig() (*Config, error) {
	root, err := utils.FindModuleRoot()
	if err != nil {
		return nil, fmt.Errorf("error encontrando la raíz del módulo: %w", err)
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	var configFile string
	switch env {
	case "railway":
		configFile = "railway-config.json"
	default:
		configFile = "local-config.json"
	}

	configPath := filepath.Join(root, "config", configFile)

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("Archivo %s no encontrado, usando local-config.json", configFile)
		configPath = filepath.Join(root, "config", "local-config.json")
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo de configuración %s: %w", configPath, err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("error parseando configuración JSON: %w", err)
	}

	log.Printf("Configuración cargada desde: %s (entorno: %s)", configFile, config.Environment)

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuración inválida: %w", err)
	}

	return &config, nil
}

func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Error fatal cargando configuración: ", err)
	}
	return config
}

func (c *Config) Validate() error {
	if c.Database.Driver == "" {
		return fmt.Errorf("driver de base de datos es requerido")
	}

	if c.Database.URL != "" {
		return nil
	}

	if c.Database.Host == "" {
		return fmt.Errorf("host de base de datos es requerido")
	}
	if c.Database.Port == "" {
		return fmt.Errorf("puerto de base de datos es requerido")
	}
	if c.Database.Username == "" {
		return fmt.Errorf("usuario de base de datos es requerido")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("contraseña de base de datos es requerida")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("nombre de base de datos es requerido")
	}

	return nil
}

func (c *Config) GetMySQLDSN() string {

	if c.Database.URL != "" {
		return c.Database.URL
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
	)

	if c.Database.SSL != "" {
		dsn += "&tls=" + c.Database.SSL
	}

	return dsn
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "railway"
}

func (c *Config) PrintConfig() {
	log.Println("=== Configuración Cargada ===")
	log.Printf("Entorno: %s", c.Environment)
	log.Printf("Servidor: %s", c.GetServerAddress())
	log.Printf("Base de Datos Driver: %s", c.Database.Driver)

	if c.Database.URL != "" {
		log.Printf("Base de Datos: URL completa configurada")
	} else {
		log.Printf("Base de Datos Host: %s:%s", c.Database.Host, c.Database.Port)
		log.Printf("Base de Datos Nombre: %s", c.Database.Name)
		log.Printf("Base de Datos Usuario: %s", c.Database.Username)
	}

	if c.Database.SSL != "" {
		log.Printf("SSL: %s", c.Database.SSL)
	}
	log.Println("=============================")
}
