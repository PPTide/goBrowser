package data

import (
	_ "embed"
	"encoding/json"
)

//go:embed entities.json
var entitiesJSON []byte

type Entity struct {
	Characters string
	Codepoints []float64
}

var (
	entitiesTmp map[string]map[string]interface{}
	entities    map[string]Entity
)

func init() {
	err := json.Unmarshal(entitiesJSON, &entitiesTmp)
	checkErr(err)
	entities = make(map[string]Entity, 0)
	//Transform map[string]map[string]interface{} into map[string]entity
	for key1, entityTmp := range entitiesTmp {
		e := Entity{}
		for key, value := range entityTmp {
			if key == "characters" {
				e.Characters = value.(string)
				continue
			}
			if key == "codepoints" {
				codepointsTmp := value.([]interface{})
				e.Codepoints = make([]float64, 0)
				for _, codepoint := range codepointsTmp {
					e.Codepoints = append(e.Codepoints, codepoint.(float64))
				}
			}
		}
		entities[key1] = e
	}
}

func GetEntities() map[string]Entity {
	return entities
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
