package config

type Config struct {
	Port string `env:"PORT"`

	MusicInfoURL string `env:"MUSIC_INFO_URL"`

	LogLevel int `env:"LOG_LEVEL"`

	PgDSN         string `env:"SONG_LIBRARY_PG_DSN"`
	PgMaxOpenConn int    `env:"PG_MAX_OPEN_CONN"`
}
