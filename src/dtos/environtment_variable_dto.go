package dtos

type EnvirontmentVariable struct {
	ApiEnv    string
	DBName    string
	DBUser    string
	DBPort    string
	DBHost    string
	DBPass    string
	DBSSL     string
	AppHost   string
	AppPort   string
	JwtSecret string
}
