package handle

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/chand1012/hooky/pkg/config"
	"github.com/charmbracelet/log"
)

type OptionMap = map[string]*discordgo.ApplicationCommandInteractionDataOption

func ParseOptions(options []*discordgo.ApplicationCommandInteractionDataOption) (om OptionMap) {
	om = make(OptionMap)
	for _, opt := range options {
		om[opt.Name] = opt
	}
	return
}

func Command(commands config.Commands, s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	opts := ParseOptions(data.Options)
	command := commands.FindByName(data.Name)
	if command == nil {
		log.Errorf("could not find command: %s", data.Name)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not find command",
			},
		})
		if err != nil {
			log.Errorf("could not respond to interaction: %s", err)
		}
		return
	}

	log.Infof("handling command: %s", command.Name)

	// construct the body
	body := make(map[string]interface{})
	for _, param := range command.Body {
		if _, ok := opts[param.Name]; ok {
			body[param.Name] = opts[param.Name].Value
		}
	}
	var bodyBuffer *bytes.Buffer
	if len(body) > 0 && command.BodyTemplate == "" {
		bodyBuffer = new(bytes.Buffer)
		if err := json.NewEncoder(bodyBuffer).Encode(body); err != nil {
			log.Errorf("could not encode body: %s", err)
			return
		}
	} else if len(body) > 0 {
		bodyBuffer = new(bytes.Buffer)
		tmpl, err := template.New("json").Parse(command.BodyTemplate)
		if err != nil {
			log.Errorf("could not parse body template: %s", err)
			return
		}
		if err := tmpl.Execute(bodyBuffer, body); err != nil {
			log.Errorf("could not execute body template: %s", err)
			return
		}
	}

	// construct the query params
	query := make(map[string]string)
	for _, param := range command.Query {
		query[param.Name] = opts[param.Name].StringValue()
	}

	// TODO add form data here

	// construct the request
	req, err := http.NewRequest(command.Method, command.URL, bodyBuffer)
	if err != nil {
		log.Errorf("could not create request: %s", err)
		return
	}
	for key, value := range command.Headers {
		req.Header.Add(key, value)
	}
	// add the query params
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("could not send request: %s", err)
		return
	}

	defer resp.Body.Close()

	// read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("could not read response body: %s", err)
		return
	}

	if command.ParseJSON != nil {
		respData, err := command.ParseJSON.Parse(respBody)
		if err != nil {
			log.Errorf("could not parse response: %s", err)
			return
		}
		respJSON, err := json.Marshal(respData)
		if err != nil {
			log.Errorf("could not marshal response: %s", err)
			return
		}
		respBody = respJSON
		if command.ResponseTemplate != "" {
			tmpl, err := template.New("json").Parse(command.ResponseTemplate)
			if err != nil {
				log.Errorf("could not parse response template: %s", err)
				return
			}
			respBuffer := new(bytes.Buffer)
			if err := tmpl.Execute(respBuffer, respData); err != nil {
				log.Errorf("could not execute response template: %s", err)
				return
			}
			respBody = respBuffer.Bytes()
		}
	}

	// first, get the response body as a string
	respString := string(respBody)
	// attempt to parse the response body as JSON
	// if it fails, just send the response body as a string
	var respJSON map[string]interface{}
	if err := json.Unmarshal(respBody, &respJSON); err != nil {
		log.Debugf("could not parse response body as JSON: %s", err)
		respJSON = nil
	}

	if respJSON != nil {
		// if the response is not nil, send the response as pretty JSON
		prettyResp, err := json.MarshalIndent(respJSON, "", "  ")
		if err == nil {
			respString = "```json\n" + string(prettyResp) + "\n```"
		} else {
			log.Debugf("could not pretty print response: %s", err)
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: respString,
		},
	})

	if err != nil {
		log.Errorf("could not respond to interaction: %s", err)
	}
}
