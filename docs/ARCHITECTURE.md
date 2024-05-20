# Architecture
## Configuration

Hooky is configured using a JSON file. The configuration is a directory of JSON files, one for each command. Each command file contains the configuration for that command.

Here is an example configuration for a simple webhook command:
```json
{
  "name": "example",
  "description": "An example command",
  "method": "POST",
  "url": "https://example.com/webhook",
  "body": [{ // can also be "query" for query parameters
    "name": "content",
    "description": "The content of the message",
    "type": "string",
    "required": true // optional. Default is false
  }],
  "headers": {
    "Authorization": "Bearer example-token",
    "Content-Type": "application/json"
  }
  // TODO: Add support for response handling
  // TODO: Add support for file uploads
  // TODO: Add support for multi-level complex JSON objects
}
```
The name of the file within the config directory does not matter, they will all be loaded.
