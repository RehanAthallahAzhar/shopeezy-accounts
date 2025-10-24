package configs

type RedisConfig struct {
	Password string `env:"REDIS_PASSWORD"`
	DB       int    `env:"REDIS_DB,required"`
}
