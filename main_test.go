package main

import (
	"testing"
)

func TestGenerateExport_WithDPT(t *testing.T) {
	input, objTypes := getTestContext()

	export, err := GenerateExport(input, objTypes, true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	lightGroup := export.GroupRanges[3].GroupRanges[0]
	// Check the second object (offset 10)
	if lightGroup.GroupAddresses[2].DPT != "DPST-1-1" {
		t.Errorf("Expected DPT DPST-1-1, got %s", lightGroup.GroupAddresses[2].DPT)
	}
}

func TestGenerateExport_WithoutDPT(t *testing.T) {
	input, objTypes := getTestContext()

	export, err := GenerateExport(input, objTypes, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	lightGroup := export.GroupRanges[3].GroupRanges[0]
	if lightGroup.GroupAddresses[2].DPT != "" {
		t.Errorf("Expected empty DPT, got %s", lightGroup.GroupAddresses[2].DPT)
	}
}

func TestGenerateExport_MissingType(t *testing.T) {
	input := Input{
		Trades: []Trade{
			{ID: "light", Name: "Licht"},
		},
		Rooms: []map[string]interface{}{
			{
				"name": "Wohnzimmer",
				"light": []interface{}{
					map[string]interface{}{
						"name": "Decke",
						"type": "unknown",
					},
				},
			},
		},
	}

	objTypes := ObjectTypes{
		"light": {
			"dimmable": ObjectTypeDefinition{
				ReservedAddresses: 10,
				Functions:         []FunctionDefinition{{Name: "schalten"}},
			},
		},
	}

	_, err := GenerateExport(input, objTypes, false)
	if err == nil {
		t.Error("Expected error for missing type, got nil")
	}
}

func getTestContext() (Input, ObjectTypes) {
	input := Input{
		Trades: []Trade{
			{ID: "light", Name: "Licht"},
		},
		Rooms: []map[string]interface{}{
			{
				"name": "Wohnzimmer",
				"light": []interface{}{
					map[string]interface{}{
						"name": "Decke",
						"type": "dimmable",
					},
					map[string]interface{}{
						"name": "Wand",
						"type": "dimmable",
					},
				},
			},
		},
	}

	objTypes := ObjectTypes{
		"light": {
			"dimmable": ObjectTypeDefinition{
				ReservedAddresses: 10,
				Functions: []FunctionDefinition{
					{Name: "schalten", DPT: "DPST-1-1"},
					{Name: "dimmen", DPT: "DPST-3-7"},
				},
			},
		},
	}
	return input, objTypes
}
