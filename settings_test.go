package main

import (
	"github.com/mailru/easyjson"
	"testing"
)

func TestParsingSettingsWithNoValueProvided(t *testing.T) {
	rawSettings := []byte(`{}`)
	settings := &Settings{}
	if err := easyjson.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if len(settings.PropagatedLabels) != 0 {
		t.Errorf("PropagatedLabels should contains zero labels after unmarshal")
	}

	valid, _ := settings.Valid()
	if valid {
		t.Errorf("Empty settings should not be valid")
	}
}

func TestParsingSettingsWithEmptyStringLabel(t *testing.T) {
	rawSettings := []byte(`{"propagatedLabels": ["label", ""]}`)
	settings := &Settings{}
	if err := easyjson.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if len(settings.PropagatedLabels) != 2 {
		t.Errorf("PropagatedLabels should contains two labels after unmarshal")
	}

	valid, _ := settings.Valid()
	if valid {
		t.Errorf("Empty label string should not be valid")
	}
}

func TestParsingSettingsWithNoLabels(t *testing.T) {
	rawSettings := []byte(`{"propagatedLabels": []}`)
	settings := &Settings{}
	if err := easyjson.Unmarshal(rawSettings, settings); err != nil {
		t.Errorf("Unexpected error %+v", err)
	}

	if len(settings.PropagatedLabels) != 0 {
		t.Errorf("PropagatedLabels should contains zero labels after unmarshal")
	}

	valid, _ := settings.Valid()
	if valid {
		t.Errorf("At least one label must be provided")
	}
}
