package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

type Gewerk struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Object struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Input struct {
	Gewerke []Gewerk                 `json:"gewerke"`
	Raeume  []map[string]interface{} `json:"raeume"`
}

type ObjectTypes map[string]map[string][]string

// XML structures
type GroupAddressExport struct {
	XMLName     xml.Name     `xml:"GroupAddress-Export"`
	XMLNS       string       `xml:"xmlns,attr"`
	GroupRanges []GroupRange `xml:"GroupRange"`
}

type GroupRange struct {
	Name           string         `xml:"Name,attr"`
	RangeStart     int            `xml:"RangeStart,attr"`
	RangeEnd       int            `xml:"RangeEnd,attr"`
	GroupRanges    []GroupRange   `xml:"GroupRange,omitempty"`
	GroupAddresses []GroupAddress `xml:"GroupAddress,omitempty"`
}

type GroupAddress struct {
	Name    string `xml:"Name,attr"`
	Address string `xml:"Address,attr"`
}

func main() {
	// 1. Load knx-object-types.json
	typesData, err := os.ReadFile("knx-object-types.json")
	if err != nil {
		log.Fatalf("Error reading knx-object-types.json: %v", err)
	}
	var objTypes ObjectTypes
	if err := json.Unmarshal(typesData, &objTypes); err != nil {
		log.Fatalf("Error unmarshaling knx-object-types.json: %v", err)
	}

	// 2. Load input.json
	inputData, err := os.ReadFile("input.json")
	if err != nil {
		log.Fatalf("Error reading input.json: %v", err)
	}
	var input Input
	if err := json.Unmarshal(inputData, &input); err != nil {
		log.Fatalf("Error unmarshaling input.json: %v", err)
	}

	// 3. Process
	export, err := GenerateExport(input, objTypes)
	if err != nil {
		log.Fatalf("Error generating export: %v", err)
	}

	// 4. Output XML
	output, err := xml.MarshalIndent(export, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling XML: %v", err)
	}

	header := []byte(xml.Header)
	finalOutput := append(header, output...)

	err = os.WriteFile("output.xml", finalOutput, 0644)
	if err != nil {
		log.Fatalf("Error writing output.xml: %v", err)
	}
}

func GenerateExport(input Input, objTypes ObjectTypes) (GroupAddressExport, error) {
	export := GroupAddressExport{
		XMLNS: "http://knx.org/xml/ga-export/01",
	}

	mainGroupIndex := 0

	// Add central groups
	centrals := []string{"Zentral", "UG-Zentral", "EG-Zentral"}
	for _, name := range centrals {
		export.GroupRanges = append(export.GroupRanges, createMainGroup(name, mainGroupIndex))
		mainGroupIndex++
	}

	// Process rooms
	for _, roomMap := range input.Raeume {
		roomName, ok := roomMap["name"].(string)
		if !ok {
			continue
		}

		mainGroup := createMainGroup(roomName, mainGroupIndex)

		middleGroupIndex := 0
		for _, gewerk := range input.Gewerke {
			objectsRaw, ok := roomMap[gewerk.ID]
			if !ok {
				continue
			}

			// Parse objects
			objectsData, _ := json.Marshal(objectsRaw)
			var objects []Object
			json.Unmarshal(objectsData, &objects)

			if len(objects) == 0 {
				continue
			}

			middleGroup := createMiddleGroup(gewerk.Name, mainGroupIndex, middleGroupIndex)
			subGroupIndex := 0

			for _, obj := range objects {
				// Get functions for this type
				functions, found := objTypes[gewerk.ID][obj.Type]
				if !found {
					return GroupAddressExport{}, fmt.Errorf("type '%s' for gewerk '%s' not found in knx-object-types.json", obj.Type, gewerk.ID)
				}

				for _, fn := range functions {
					if subGroupIndex > 255 {
						return GroupAddressExport{}, fmt.Errorf("too many group addresses in middle group %s/%s", roomName, gewerk.Name)
					}

					gaName := fmt.Sprintf("%s_%s_%s-%s", roomName, gewerk.Name, obj.Name, fn)
					gaAddress := fmt.Sprintf("%d/%d/%d", mainGroupIndex, middleGroupIndex, subGroupIndex)

					middleGroup.GroupAddresses = append(middleGroup.GroupAddresses, GroupAddress{
						Name:    gaName,
						Address: gaAddress,
					})
					subGroupIndex++
				}
			}

			mainGroup.GroupRanges = append(mainGroup.GroupRanges, middleGroup)
			middleGroupIndex++
		}

		export.GroupRanges = append(export.GroupRanges, mainGroup)
		mainGroupIndex++
	}

	return export, nil
}

func createMainGroup(name string, index int) GroupRange {
	start := index * 2048
	if index == 0 {
		start = 1 // 0/0/0 is reserved
	}
	return GroupRange{
		Name:       name,
		RangeStart: start,
		RangeEnd:   (index+1)*2048 - 1,
	}
}

func createMiddleGroup(name string, mainIndex, middleIndex int) GroupRange {
	start := mainIndex*2048 + middleIndex*256
	return GroupRange{
		Name:       name,
		RangeStart: start,
		RangeEnd:   mainIndex*2048 + (middleIndex+1)*256 - 1,
	}
}
