package database

// Configuracion representa la configuración completa del archivo YAML
type Configuracion struct {
	Database struct {
		Principal  ConfiguracionDB `yaml:"principal"`
		Secundaria ConfiguracionDB `yaml:"secundaria"`
	} `yaml:"database"`
}

// ConfiguracionDB representa la configuración de una base de datos específica
type ConfiguracionDB struct {
	Host            string            `yaml:"host"`
	Port            int               `yaml:"port"`
	Name            string            `yaml:"name"`
	User            string            `yaml:"user"`
	Password        string            `yaml:"password"`
	Encrypt         bool              `yaml:"encrypt"`
	TrustServerCert bool              `yaml:"trust_server_certificate"`
	Pool            ConfiguracionPool `yaml:"pool"`
}

// ConfiguracionPool representa la configuración del pool de conexiones
type ConfiguracionPool struct {
	Min                 int `yaml:"min"`
	Max                 int `yaml:"max"`
	IdleTimeoutMs       int `yaml:"idle_timeout_ms"`
	ConnectionTimeoutMs int `yaml:"connection_timeout_ms"`
}
