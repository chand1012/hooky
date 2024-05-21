set dotenv-load

GUILD_ID := env_var("GUILD_ID")
APP_ID := env_var("APP_ID")
BOT_TOKEN := env_var("BOT_TOKEN")

dev:
  go run main.go -guild {{GUILD_ID}} -app {{APP_ID}} -token {{BOT_TOKEN}}

build:
  go build -o bot main.go

tidy:
  go mod tidy

clean:
  rm -f bot
  go clean -cache -testcache
