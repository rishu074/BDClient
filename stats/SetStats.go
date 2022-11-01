package stats

import (
	"encoding/json"
	"os"
)

func SetStats(s bool) {
	type Data struct {
		Running bool
	}

	var data Data
	if s {
		data.Running = true
	} else {
		data.Running = false
	}

	b, _ := json.MarshalIndent(data, "", "")
	os.WriteFile("./stats/stats.json", b, 0644)
}
