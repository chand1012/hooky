# Hooky - configurable REST Discord bot

Hooky is a bot that allows you to configure any webhook or REST request into a Discord command that can be triggered by users. It is designed to be easy to use and flexible, allowing you to configure any webhook or REST request you want.

## Why'd you make this?

I make a lot of projects. There's usually one thing those projects all have in common - a REST API. I wanted a way to quickly use my REST APIs from some sort of interface without having to write either a new Discord bot or a new web interface for each project. I also wanted to be able to configure these commands without having to write any code (JSON isn't code, right?).

This means that in theory, if you just write a REST API, you can use Hooky to interact with it or even turn it into its own standalone Discord bot.

For example, something I use a lot is self-hosted no-code tools like [n8n](https://n8n.io) and [Dify](https://github.com/langgenius/dify). Both offer no-code options to interact with, however I wanted a single interface to interact with them, as well as any other REST interface. Hooky allows me to do that.

## Features

Checked means it's implemented, unchecked means it's planned.

- [x] Configurable commands
- [x] JSON parsing and requests
- [ ] JSON list parsing
- [x] Request templates
- [x] Response templates
- [x] Query params
- [ ] Form data
- [ ] File uploads
- [ ] Rate limiting
- [ ] URL parameters

## Setup

### Prerequisites
- Go 1.22 or later or Docker
- A Discord bot token
- A Discord app ID
- (Optional) A Discord guild ID
- (Optional) Just

```bash
git clone https://github.com/chand1012/hooky.git
cd hooky
```

### Building
```bash
just build
## or if you don't have Just
go build -o hooky main.go
## or if you want to use Docker
docker build -t hooky .
```

### Running with Go
```bash
# guild and config-dir are optional. Token, App, and Guild can be specified via flags or environment variables.
./hooky -token <bot-token> -app <app-id> -config-dir <config-dir> -guild <guild-id>
```

If you have Just installed, you can also run `just start` to start the bot. This will also allow you to load a `.env` file with the following variables:
```
BOT_TOKEN=<bot-token>
APP_ID=<app-id>
GUILD_ID=<guild-id>
```

```bash
# load .env, build, and run
just start
```

### Running with Docker
```bash
docker run -e BOT_TOKEN=<bot-token> -e APP_ID=<app-id> -e GUILD_ID=<guild-id> -v /path/to/config:/app/config hooky
```

## Configuration

For the Guild ID, bot token, and app ID, 

Hooky uses a configuration file in JSON format to define commands and their corresponding REST requests. The configuration file is loaded from the directory specified by the `--config-dir` flag. All files in the directory are loaded and parsed as command configurations, name doesn't matter. One file is a single command configuration.

### Command Configuration

A command configuration defines a single command that can be triggered by a user. A command configuration is a JSON object with the following properties:

* `name`: A string that specifies the name of the command. This is the command that users will trigger to execute the corresponding REST request.
* `description`: A string that provides a brief description of the command.
* `method`: A string that specifies the HTTP method to use for the REST request (e.g. `GET`, `POST`, `PUT`, `DELETE`, etc.).
* `url`: A string that specifies the URL of the REST request.
* `body`: An array of `Param` objects that specify the parameters to include in the request body.
* `body_template`: A string that specifies a template to use for generating the request body. The template can use variables from the `body` array. Should be in the form of your request body, with variables in the form of `{{.variable_name}}`.
* `query`: An array of `Param` objects that specify the query parameters to include in the request.
* `form`: An array of `Param` objects that specify the form data to include in the request.
* `headers`: An object that specifies the headers to include in the request.
* `parse_json`: An object that specifies how to parse the response JSON.
* `response_template`: A string that specifies a template to use for generating the response. The template can use variables from the parsed JSON response. Should be in the form of your response, with variables in the form of `{{.variable_name}}`.

### Param Configuration

A `Param` object specifies a single parameter for a command. A `Param` object has the following properties:

* `name`: A string that specifies the name of the parameter.
* `description`: A string that provides a brief description of the parameter.
* `type`: A string that specifies the type of the parameter (e.g. `string`, `integer`, `boolean`, etc.).
* `required`: A boolean that specifies whether the parameter is required.
* `default`: A value that specifies the default value of the parameter.
* `options`: An array of values that specifies the valid options for the parameter.

### Example Configuration

Here is an example configuration file that defines a single command:
```json
{
  "name": "example",
  "description": "This is a test command",
  "method": "POST",
  "url": "https://httpbin.org/anything",
  "body": [
    {
      "name": "content",
      "description": "This is a test command",
      "type": "string",
      "required": true
    },
    {
      "name": "bees",
      "description": "save the bees",
      "required": false,
      "type": "boolean",
      "default": false
    },
    {
      "name": "select",
      "description": "This is a selection",
      "required": false,
      "type": "string",
      "options": ["thing 1", "thing 2"],
      "default": "thing 2"
    }
  ],
  "body_template": "{ \"content\": \"{{ .content }}\", \"dummy\": \"am dummy\", \"bees\": \"{{ .bees }}\", \"select\": \"{{.select}}\" }",
  "parse_json": {
    "response_content": ".json.content",
    "bees": ".json.bees",
    "select": ".json.select"
  },
  "response_template": "{{.response_content}}\n\nBees? {{.bees}}\n\nSelection: {{.select}}"
}

```
This configuration defines a command called `example` that makes a `POST` request to `https://httpbin.org/anything`. The request body includes two parameters: `content` and `bees`. The `content` parameter is required and has a type of `string`, while the `bees` parameter is optional and has a type of `boolean`. The request also includes a query parameter `api_key` that is required. The response is parsed as JSON and the `response_template` is used to generate the final response.
