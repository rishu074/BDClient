package stats

import (
	"encoding/json"
	"os"
)

func GetStats() (bool, error) {
	file, err := os.Open("./stats/stats.json")
	if err != nil {
		return false, err
	}

	defer file.Close()
	decoder := json.NewDecoder(file)

	var data map[string]interface{}
	decoder.Decode(&data)

	return data["Running"].(bool), nil
}
