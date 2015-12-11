package config

import "testing"

func checkEqualString(t *testing.T, testName, got, want string) {
	if got != want {
		t.Errorf("Test: %s failed got %s, wanted: %s", testName, got, want)
	}
}

func checkEqualInt(t *testing.T, testName string, got, want int) {
	if got != want {
		t.Errorf("Test: %s failed got %d, wanted: %d", testName, got, want)
	}
}

func TestOuputFormats(t *testing.T) {
	output := ParseConfigOutputFormat("test/output_config_1.yml")
	checkEqualString(t, "type-0", output.Outputs[0].Filetype, "jpeg")
	checkEqualString(t, "type-1", output.Outputs[1].Filetype, "text")
	checkEqualString(t, "filename-0", output.Outputs[0].Filename, "output.jpg")
	checkEqualString(t, "filename-1", output.Outputs[1].Filename, "output.txt")
	checkEqualInt(t, "width-0", output.Outputs[0].Width, 1024)
	checkEqualInt(t, "filter-0", output.Outputs[2].Filter.Latitude.Min, 576906339)
}

func TestParseConfigInput(t *testing.T) {
	inputs := ParseConfigInput("test/input_config_1.yml")
	checkEqualString(t, "type-0", inputs.Inputs[0].Type, "json")
	checkEqualString(t, "type-1", inputs.Inputs[1].Type, "csv")
	checkEqualString(t, "type-2", inputs.Inputs[2].Type, "gpx")

}
