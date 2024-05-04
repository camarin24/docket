package adapters

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Tika struct {
	serverEndpoint string
}

type TikaConfig struct {
	ServerEndpoint string
}

func (t *Tika) GetFileMetadata(path string) (*map[string]interface{}, *string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	baseURL, err := url.JoinPath(t.serverEndpoint, "rmeta", "text")
	if err != nil {
		return nil, nil, err
	}
	req, err := http.NewRequest("PUT", baseURL, bytes.NewBuffer(file))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, nil, err
		}

		var response []map[string]interface{}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, nil, err
		}

		content := new(string)
		parsedResponse := make(map[string]interface{})
		for k, v := range response[0] {
			if k == "X-TIKA:content" {
				bytes, err := json.Marshal(v)
				if err != nil {
					return nil, nil, err
				}
				*content = string(bytes)
				*content = strings.Replace(*content, "\\n\\n", "", -1)
				continue
			}

			// TODO: Remove all XTIKA properties
			if k != "pdf:unmappedUnicodeCharsPerPage" && k != "X-TIKA:Parsed-By-Full-Set" && k != "pdf:charsPerPage" && k != "X-TIKA:Parsed-By" {
				name := strings.Replace(k, ":", "_", -1)
				name = strings.Replace(name, ".", "_", -1)
				name = strings.Replace(name, "-", "_", -1)
				parts := strings.Split(name, "_")
				newName := ""
				for _, part := range parts {
					newName += strings.Title(part)
				}
				parsedResponse[newName] = v
			}

		}

		// TODO: Deal with multi page-metadata response (.pptx)
		return &parsedResponse, content, nil
	}

	return nil, nil, nil
}

func NewTika(config TikaConfig) *Tika {
	return &Tika{
		serverEndpoint: config.ServerEndpoint,
	}
}
