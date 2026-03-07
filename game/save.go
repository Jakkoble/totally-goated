package game

import (
	"encoding/json"
	"os"
)

const saveFile = "save.json"

type SaveData struct {
	BestScore float64 `json:"best_score"`
}

func loadSave() SaveData {
	data, err := os.ReadFile(saveFile)
	if err != nil {
		return SaveData{}
	}
	var s SaveData
	json.Unmarshal(data, &s)
	return s
}

func saveBest(score float64) {
	s := SaveData{BestScore: score}
	data, _ := json.Marshal(s)
	os.WriteFile(saveFile, data, 0644)
}
