package env

type environment struct {
	Port     int  `env:"PORT" default:"8080"`
	Debug    bool `env:"DEBUG" default:"false"`
	Database struct {
		DSN string `env:"DB_DSN"`
	}
	JWTSecret string `env:"JWT_SECRET"`
}
