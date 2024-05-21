# Architecture
## Configuration

Hooky is configured using a JSON file. The configuration is a directory of JSON files, one for each command. Each command file contains the configuration for that command.

Here is an example configuration for a simple webhook command:
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
      "type": "boolean"
    }
  ],
  "body_template": "{ \"content\": \"{{ .content }}\", \"dummy\": \"am dummy\", \"bees\": \"{{ .bees }}\" }",
  "parse_json": {
    "response_content": ".json.content",
    "bees": ".json.bees"
  },
  "response_template": "{{.response_content}}\n\nBees? {{.bees}}"
}

```
The name of the file within the config directory does not matter, they will all be loaded.
