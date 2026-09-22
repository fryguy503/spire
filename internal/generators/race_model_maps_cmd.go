package generators

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type RaceModelMapsCommand struct {
	command *cobra.Command
}

func (c *RaceModelMapsCommand) Command() *cobra.Command {
	return c.command
}

type Zone struct {
	ID         int    `json:"id"`
	LongName   string `json:"long_name"`
	ShortName  string `json:"short_name"`
	ImportType string `json:"import_type"`
}

type Model struct {
	Code           string `json:"code"`
	Gender         int    `json:"gender"`
	MinTexture     int    `json:"min_texture"`
	MaxTexture     int    `json:"max_texture"`
	MinHelmTexture int    `json:"min_helm_texture"`
	MaxHelmTexture int    `json:"max_helm_texture"`
	MinHair        int    `json:"min_hair"`
	MaxHair        int    `json:"max_hair"`
	MinBeards      int    `json:"min_beards"`
	MaxBeards      int    `json:"max_beards"`
}

type Source struct {
	SourceFile           string  `json:"source_file"`
	LoadedViaDescription string  `json:"loaded_via_description"`
	IsGlobal             bool    `json:"is_global"`
	Models               []Model `json:"models"`
	Zones                []Zone  `json:"zones"`
}

type RaceEntry struct {
	RaceId                 int      `json:"race_id"`
	Description            string   `json:"description"`
	IsPlayable             bool     `json:"is_playable"`
	MaleGenderModelCode    string   `json:"male_gender_model_code"`
	FemaleGenderModelCode  string   `json:"female_gender_model_code"`
	NeutralGenderModelCode string   `json:"neutral_gender_model_code"`
	Sources                []Source `json:"sources"`
}

type RaceData struct {
	SourceURL string      `json:"source_url"`
	Client    string      `json:"client"`
	Race      []RaceEntry `json:"races"`
}

func NewRaceModelMapsCommand() *RaceModelMapsCommand {
	i := &RaceModelMapsCommand{
		command: &cobra.Command{
			Use:   "make:race-model-maps",
			Short: "Generates race model maps from the Clumsy’s World RoF2 reference",
		},
	}

	i.command.Args = i.Validate
	i.command.RunE = i.Handle
	i.command.Flags().String("source-file", "", "Read a saved reference HTML file instead of downloading it")

	return i
}

// Handle implementation of the Command interface
const raceInventorySourceURL = "https://races.clumsysworld.com/"

func (c *RaceModelMapsCommand) Handle(cmd *cobra.Command, _ []string) error {
	sourceFile, _ := cmd.Flags().GetString("source-file")
	var contents []byte
	var err error
	if sourceFile != "" {
		contents, err = os.ReadFile(sourceFile)
	} else {
		client := &http.Client{Timeout: 30 * time.Second}
		var response *http.Response
		response, err = client.Get(raceInventorySourceURL)
		if err == nil {
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("race reference returned HTTP %d", response.StatusCode)
			}
			contents, err = io.ReadAll(io.LimitReader(response.Body, 10*1024*1024))
		}
	}
	if err != nil {
		return err
	}
	rd := c.Parse(string(contents))
	if len(rd.Race) == 0 {
		return fmt.Errorf("reference contained no race entries; existing inventory was not changed")
	}
	encoded, err := json.Marshal(rd)
	if err != nil {
		return err
	}
	file := "internal/http/staticmaps/race-inventory-map.json"
	if err := os.WriteFile(file, encoded, 0644); err != nil {
		return err
	}
	fmt.Printf("Wrote %d races to [%s]\n", len(rd.Race), file)
	return nil
}

func (c *RaceModelMapsCommand) Parse(contents string) RaceData {
	rd := RaceData{SourceURL: raceInventorySourceURL, Client: "Rain of Fear 2 (RoF2)"}
	doc := soup.HTMLParse(contents)
	for _, section := range doc.FindAll("span") {
		id := section.Attrs()["id"]
		if !strings.HasPrefix(id, "Race") {
			continue
		}
		i, err := strconv.Atoi(strings.TrimPrefix(id, "Race"))
		if err != nil {
			continue
		}

		// add New race entry
		raceEntry := RaceEntry{}
		raceEntry.RaceId = i

		// parse race section
		if section.Error == nil {

			// vars
			raceDescription := ""
			isPlayable := false

			// parse header
			header := section.Find("h3")
			if header.Error == nil {
				headerSplit1 := strings.SplitN(header.Text(), " - ", 2)
				if len(headerSplit1) > 1 {
					headerSplit2 := strings.Split(headerSplit1[1], ",")
					if len(headerSplit2) > 0 {
						raceDescription = strings.TrimSpace(headerSplit2[0])
					}
				}

				// playable
				if strings.Contains(header.Text(), "Playable") {
					isPlayable = true
				}
			}

			// gender models
			maleModel := ""
			femaleModel := ""
			neutralModel := ""

			// list items (gender models)
			li := section.FindAll("li")
			if len(li) > 0 {
				for _, entry := range li {
					// text of li
					text := entry.FullText()

					// male
					if strings.Contains(text, "Male") && !strings.Contains(text, "N/A") {
						modelSplit := strings.Split(text, "=")
						if len(modelSplit) > 0 {
							maleModel = strings.TrimSpace(modelSplit[1])
						}
					}

					// female
					if strings.Contains(text, "Female") && !strings.Contains(text, "N/A") {
						modelSplit := strings.Split(text, "=")
						if len(modelSplit) > 0 {
							femaleModel = strings.TrimSpace(modelSplit[1])
						}
					}

					// neutral
					if strings.Contains(text, "Neutral") && !strings.Contains(text, "N/A") {
						modelSplit := strings.Split(text, "=")
						if len(modelSplit) > 0 {
							neutralModel = strings.TrimSpace(modelSplit[1])
						}
					}

				}
			}

			// genders
			//pp.Println("genders")
			//pp.Println(maleModel)
			//pp.Println(femaleModel)
			//pp.Println(neutralModel)

			sources := []Source{}

			// list items
			if len(li) > 0 {
				for _, entry := range li {

					// New source entry
					source := Source{}

					// text of li
					text := entry.FullText()

					// warp string to be more parse friendly
					text = strings.ReplaceAll(text, "\n", " \n")

					// via source
					source.SourceFile = c.GetStringInBetween(text, "Via SourceFile: ", " ")
					if len(source.SourceFile) == 0 {
						// 2nd try
						source.SourceFile = c.GetStringInBetween(text, "Via Source: ", " ")
					}

					// Loaded with
					source.LoadedViaDescription = c.GetStringInBetween(text, "Loaded with ", ")")
					if len(source.LoadedViaDescription) == 0 {
						// 2nd try
						source.LoadedViaDescription = c.GetStringInBetween(text, "Loaded via", ")")
					}

					// parse lines inside of li
					lines := strings.Split(text, "\n")
					for _, line := range lines {
						// Model
						if strings.Contains(line, "Model") || strings.Contains(line, "Texture") {

							// sanitize line for more consistent parsing
							newText := line + " "
							newText = strings.ReplaceAll(newText, "  ", " ")
							newText = strings.ReplaceAll(newText, ",", " ")
							newText = strings.ReplaceAll(newText, "\n\n", " ")

							// model string
							modelString := c.GetStringInBetween(newText, "Model ", " ")

							// min / max values
							minTexture, maxTexture := c.GetMinMaxValues(c.GetStringInBetween(newText, "Textures: ", " "))
							minHelmTexture, maxHelmTexture := c.GetMinMaxValues(c.GetStringInBetween(newText, "Heads: ", " "))
							minHair, maxHair := c.GetMinMaxValues(c.GetStringInBetween(newText, "Hair: ", " "))
							minBeard, maxBeard := c.GetMinMaxValues(c.GetStringInBetween(newText, "Beards: ", " "))

							// Single-gender sources omit the explicit "Model CODE" line.
							if modelString == "" {
								codes := []string{}
								for _, code := range []string{maleModel, femaleModel, neutralModel} {
									if code != "" {
										codes = append(codes, code)
									}
								}
								if len(codes) == 1 {
									modelString = codes[0]
								}
							}
							// get gender value from string value matches
							gender := 2
							if len(femaleModel) > 0 && strings.Contains(modelString, femaleModel) {
								gender = 1
							}
							if len(maleModel) > 0 && strings.Contains(modelString, maleModel) {
								gender = 0
							}

							modelCode := modelString
							// neutral
							if modelCode == "" {
								modelCode = neutralModel
							}

							// New model
							model := Model{
								Code:           modelCode,
								Gender:         gender,
								MinTexture:     minTexture,
								MaxTexture:     maxTexture,
								MinHelmTexture: minHelmTexture,
								MaxHelmTexture: maxHelmTexture,
								MinHair:        minHair,
								MaxHair:        maxHair,
								MinBeards:      minBeard,
								MaxBeards:      maxBeard,
							}

							source.Models = append(source.Models, model)
						}
					}

					// double newline split
					if strings.Contains(text, "\n \n") {
						doubleLines := strings.Split(text, "\n \n")
						for _, line := range doubleLines {
							//pp.Println("double line")

							if strings.Contains(line, "Global") {
								source.IsGlobal = true
							}
							if strings.Contains(line, "Zone") {
								for _, zone := range strings.Split(line, "\n") {
									//

									// modify string for easier parsing
									zone = zone + " "

									zoneIdString := c.GetStringInBetween(zone, "Zone # ", " ")
									zoneLongName := c.GetStringInBetween(zone, " - ", "(")
									zoneShortName := c.GetStringInBetween(zone, "(", ")")
									importType := c.GetStringInBetween(zone, "), ", " ")
									zoneId, _ := strconv.Atoi(zoneIdString)

									newZone := Zone{
										ID:         zoneId,
										LongName:   strings.TrimSpace(zoneLongName),
										ShortName:  strings.TrimSpace(zoneShortName),
										ImportType: strings.TrimSpace(importType),
									}

									if newZone.ID != 0 {
										source.Zones = append(source.Zones, newZone)
									}

									//fmt.Printf("[%v]\n", zone)
								}
							}

							//pp.Println(line)
						}
					}

					if source.SourceFile != "" {
						sources = append(sources, source)
					}

					//pp.Println(sourceDescription)
					//pp.Println(entry.FullText())
				}
			}

			// set vars to entry
			raceEntry.MaleGenderModelCode = maleModel
			raceEntry.FemaleGenderModelCode = femaleModel
			raceEntry.NeutralGenderModelCode = neutralModel
			raceEntry.Sources = sources
			raceEntry.Description = raceDescription
			raceEntry.IsPlayable = isPlayable

			// append entry
			rd.Race = append(rd.Race, raceEntry)

			//pp.Println(raceDescription)
			//pp.Println(isPlayable)
			//fmt.Println(section.HTML())
		}

	}

	return rd
}

// Validate implementation of the Command interface
func (c *RaceModelMapsCommand) Validate(_ *cobra.Command, _ []string) error {
	return nil
}

func (c *RaceModelMapsCommand) GetStringInBetween(value string, a string, b string) string {
	firstSplit := strings.Split(value, a)
	if len(firstSplit) > 1 {
		secondSplit := strings.Split(firstSplit[1], b)
		if len(secondSplit) > 0 {
			return strings.TrimSpace(secondSplit[0])
		}
	}

	return ""
}

func (c *RaceModelMapsCommand) GetMinMaxValues(between string) (int, int) {
	if strings.Contains(between, "-") {
		split := strings.Split(between, "-")
		if len(split) > 0 {
			minValue := strings.TrimSpace(split[0])
			maxValue := strings.TrimSpace(split[1])
			minInt, _ := strconv.Atoi(minValue)
			maxInt, _ := strconv.Atoi(maxValue)

			return minInt, maxInt
		}
	}

	return 0, 0
}
