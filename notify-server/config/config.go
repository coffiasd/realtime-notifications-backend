package config

type Conf struct {
	App   AppConfig
	Log   LogConfig
	Redis RedisConfig
}

type AppConfig struct {
	AppName     string
	HttpListen  string
	AlgoRandUrl string
}

type LogConfig struct {
	LogPath  string
	FileName string
}

type RedisConfig struct {
	Addr              string
	MaxIdle           int
	MaxActive         int
	Password          string
	DB                int
	ConnectionTimeout int
	ReadTimeout       int
	WriteTimeout      int
	IdleTimeout       int
}