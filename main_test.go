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
				},
			},
		},
	}

	objTypes := ObjectTypes{
		"light": {
			"dimmable": []string{"schalten", "dimmen"},
		},
	}

	export, err := GenerateExport(input, objTypes)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// 3 centrals + 1 room = 4 main groups
	if len(export.GroupRanges) != 4 {
		t.Errorf("Expected 4 main groups, got %d", len(export.GroupRanges))
	}

	roomGroup := export.GroupRanges[3]
	if roomGroup.Name != "Wohnzimmer" {
		t.Errorf("Expected room name 'Wohnzimmer', got %s", roomGroup.Name)
	}

	if len(roomGroup.GroupRanges) != 1 {
		t.Errorf("Expected 1 middle group, got %d", len(roomGroup.GroupRanges))
	}

	lightGroup := roomGroup.GroupRanges[0]
	if lightGroup.Name != "Licht" {
		t.Errorf("Expected middle group name 'Licht', got %s", lightGroup.Name)
	}

	if len(lightGroup.GroupAddresses) != 2 {
		t.Errorf("Expected 2 group addresses, got %d", len(lightGroup.GroupAddresses))
	}

	expectedGA := "Wohnzimmer_Licht_Decke-schalten"
	if lightGroup.GroupAddresses[0].Name != expectedGA {
		t.Errorf("Expected GA name %s, got %s", expectedGA, lightGroup.GroupAddresses[0].Name)
	}

	expectedAddr := "3/0/0"
	if lightGroup.GroupAddresses[0].Address != expectedAddr {
		t.Errorf("Expected address %s, got %s", expectedAddr, lightGroup.GroupAddresses[0].Address)
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
			"dimmable": []string{"schalten"},
		},
	}

	_, err := GenerateExport(input, objTypes)
	if err == nil {
		t.Error("Expected error for missing type, got nil")
	}
}
