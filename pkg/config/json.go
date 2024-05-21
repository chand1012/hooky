package config

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/itchyny/gojq"
)

type ParseJSONResponse map[string]string

func (p *ParseJSONResponse) Parse(body []byte) (map[string]string, error) {
	respData := make(map[string]string)
	var input map[string]interface{}
	if err := json.Unmarshal(body, &input); err != nil {
		return nil, err
	}
	for key, value := range *p {
		query, err := gojq.Parse(value)
		if err != nil {
			log.Errorf("could not parse jq query: %s", err)
			continue
		}
		iter := query.Run(input)
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
					break
				}
				log.Errorf("could not run jq query: %s", err)
			}
			respData[key] = fmt.Sprintf("%v", v)
		}
	}

	return respData, nil
}
