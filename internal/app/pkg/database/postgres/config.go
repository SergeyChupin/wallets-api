package postgres

type Config struct {
	Url string `mapstructure:"url"`
}

func NewConfig() Config {
	return Config{
		Url: "postgres://127.0.0.1:5432/postgres",
	}
}
