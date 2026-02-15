package main

import (
	"testing"
)

func TestGenerateExport_Success(t *testing.T) {
	input := Input{
		Trades: []Trade{
			{ID: "light", Name: "Licht"},
		},
		Rooms: []map[string]interface{}{
			{
				"name": "Wohnzimmer",
				"light": []interface{}{
					map[string]interface{}{
						"id":   "l1",
						"name": "Decke",
						"type": "dimmable",
					},
					map[string]interface{}{
						"id":   "l2",
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

	export, err := GenerateExport(input, objTypes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	roomGroup := export.GroupRanges[3]
	lightGroup := roomGroup.GroupRanges[0]

	// Total 4 addresses (2 per object)
	if len(lightGroup.GroupAddresses) != 4 {
		t.Errorf("Expected 4 group addresses, got %d", len(lightGroup.GroupAddresses))
	}

	// First object first GA
	if lightGroup.GroupAddresses[0].Address != "3/0/0" {
		t.Errorf("Expected 3/0/0, got %s", lightGroup.GroupAddresses[0].Address)
	}

	// Second object first GA should be at offset 10
	if lightGroup.GroupAddresses[2].Address != "3/0/10" {
		t.Errorf("Expected 3/0/10, got %s", lightGroup.GroupAddresses[2].Address)
	}
	
	if lightGroup.GroupAddresses[2].DPT != "DPST-1-1" {
		t.Errorf("Expected DPT DPST-1-1, got %s", lightGroup.GroupAddresses[2].DPT)
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
						"id":   "l1",
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

	_, err := GenerateExport(input, objTypes)
	if err == nil {
		t.Error("Expected error for missing type, got nil")
	}
}
