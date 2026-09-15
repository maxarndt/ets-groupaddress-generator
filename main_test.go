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

func TestGenerateExport_WithoutObjectName(t *testing.T) {
	input := Input{
		Trades: []Trade{{ID: "clima", Name: "Heizung"}},
		Rooms: []map[string]interface{}{
			{
				"name": "EG-Bad",
				"clima": []interface{}{
					map[string]interface{}{"type": "general"},
				},
			},
		},
	}
	objTypes := ObjectTypes{
		"clima": {
			"general": {
				ReservedAddresses: 8,
				Functions:         []FunctionDefinition{{Name: "Ist-Temperatur"}},
			},
		},
	}

	export, err := GenerateExport(input, objTypes, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	got := export.GroupRanges[3].GroupRanges[0].GroupAddresses[0].Name
	want := "EG-Bad_Heizung-Ist-Temperatur"
	if got != want {
		t.Errorf("Expected group address name %q, got %q", want, got)
	}
}

func TestGenerateExport_ReservesMiddleGroupSlotForMissingTrade(t *testing.T) {
	input := Input{
		Trades: []Trade{
			{ID: "light", Name: "Licht"},
			{ID: "shading", Name: "Beschattung"},
			{ID: "clima", Name: "Klima"},
		},
		Rooms: []map[string]interface{}{
			{
				"name":  "UG-Vorratsraum",
				"light": []interface{}{map[string]interface{}{"type": "switch"}},
				"clima": []interface{}{map[string]interface{}{"type": "general"}},
			},
		},
	}
	objTypes := ObjectTypes{
		"light": {
			"switch": {Functions: []FunctionDefinition{{Name: "schalten"}}},
		},
		"clima": {
			"general": {Functions: []FunctionDefinition{{Name: "Ist-Temperatur"}}},
		},
	}

	export, err := GenerateExport(input, objTypes, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	room := export.GroupRanges[3]
	if len(room.GroupRanges) != 2 {
		t.Fatalf("Expected two middle groups, got %d", len(room.GroupRanges))
	}

	clima := room.GroupRanges[1]
	if clima.Name != "Klima" || clima.RangeStart != 3*2048+2*256 {
		t.Errorf("Expected Klima in middle-group slot 2, got name %q at RangeStart %d", clima.Name, clima.RangeStart)
	}
	if got := clima.GroupAddresses[0].Address; got != "3/2/0" {
		t.Errorf("Expected climate group address 3/2/0, got %s", got)
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
