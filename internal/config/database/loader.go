package database

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)


// CargarConfiguracion carga la configuración desde el archivo YAML
func CargarConfiguracion(rutaConfig string) (*Configuracion, error) {
	rutaConfig = os.ExpandEnv(rutaConfig)
	contenido, err := os.ReadFile(rutaConfig)
	if err != nil {
		return nil, fmt.Errorf("error al leer archivo de configuración: %w", err)
	}

	contenidoExpandido := os.ExpandEnv(string(contenido))

	var config Configuracion
	if err := yaml.Unmarshal([]byte(contenidoExpandido), &config); err != nil {
		return nil, fmt.Errorf("error al parsear YAML: %w", err)
	}

	return &config, nil
}
