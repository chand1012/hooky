set dotenv-load

GUILD_ID := env_var("GUILD_ID")
APP_ID := env_var("APP_ID")
BOT_TOKEN := env_var("BOT_TOKEN")

dev:
  go run main.go -app {{APP_ID}} -token {{BOT_TOKEN}} -guild {{GUILD_ID}}

build:
  go build -o hooky main.go

start: build
  ./hooky


tidy:
  go mod tidy

clean:
  rm -f bot
  go clean -cache -testcache

build-docker:
  docker build -t chand1012/hooky .

run-docker:
  docker run -e BOT_TOKEN={{BOT_TOKEN}} -e APP_ID={{APP_ID}} -e GUILD_ID={{GUILD_ID}} -v $(pwd)/config:/app/config chand1012/hooky
