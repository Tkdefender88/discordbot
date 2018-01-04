package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Tkdefender88/officerDva/config"
	"github.com/bwmarrin/discordgo"
)

type result struct {
	LookupList []lookup `json:"list"`
	ResultType string   `json:"result_type"`
}

type lookup struct {
	Word       string `json:"word"`
	Definition string `json:"definition"`
	Example    string `json:"example"`
	Author     string `json:"author"`
	Thumbup    int    `json:"thumbs_up"`
	Thumbdown  int    `json:"thumbs_down"`
}

const (
	apiBase    string = "http://api.urbandictionary.com"
	apiVersion string = "v0"
)

func udLookup(s *discordgo.Session, m *discordgo.MessageCreate, msgList []string) {
	//The search in the url needs to have words joined by plus signs
	search := strings.Join(msgList[1:], "+")
	url := fmt.Sprintf("%s/%s/define?term=%s", apiBase, apiVersion, search)
	//Get the json data from the web
	lookupInfo := searchUD(url)
	//Parse the json data for definition author rating and example
	res, err := parseJSONData(lookupInfo)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//If no result is found then send an error message to chat and stop.
	if res.ResultType == "no_results" {
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Color: config.EmbedColor,

			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Error",
					Value: "Definition not found",
				},
			},
		})
		return
	}

	//Send result as an embeded message
	embedUDresult(s, m, &res.LookupList[0])
}

func embedUDresult(s *discordgo.Session, m *discordgo.MessageCreate, data *lookup) {
	rating := fmt.Sprintf(":+1:`%d` :-1:`%d`", data.Thumbup, data.Thumbdown)
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Color: config.EmbedColor,

		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Definition",
				Value: data.Definition,
			},
			{
				Name:  "Example",
				Value: data.Example,
			},
			{
				Name:   "Rating",
				Value:  rating,
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  data.Author,
				Inline: true,
			},
		},
	})

}

//Sends an HTTP get request to the urban dictionary api and returns the json data that it
//receives as a response
func searchUD(url string) (body []byte) {
	data, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body, readErr := ioutil.ReadAll(data.Body)

	if readErr != nil {
		fmt.Println(readErr.Error())
		return
	}

	return body
}

func parseJSONData(data []byte) (res *result, err error) {

	jsonErr := json.Unmarshal([]byte(data), &res)

	if jsonErr != nil {
		fmt.Println(jsonErr.Error())
		return res, err
	}

	return res, nil
}
