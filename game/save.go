package game

import (
	"encoding/json"
	"os"
)

const saveFile = "save.json"

type SaveData struct {
	BestScore         int            `json:"best_score"`
	BestMeters        int            `json:"best_meters"`
	BestBells         int            `json:"best_bells"`
	TotalBells        int            `json:"total_bells"`
	UnlockedCosmetics []int          `json:"unlocked_cosmetics"`
	EquippedCosmetics map[string]int `json:"equipped_cosmetics"`
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

func saveBest(score int, meters int, bells int, totalBells int, unlockedCosmetics []int, equippedCosmetics map[string]int) {
	cur := loadSave()
	if score > cur.BestScore {
		cur.BestScore = score
	}
	if meters > cur.BestMeters {
		cur.BestMeters = meters
	}
	if bells > cur.BestBells {
		cur.BestBells = bells
	}
	cur.TotalBells = totalBells
	cur.UnlockedCosmetics = unlockedCosmetics
	cur.EquippedCosmetics = equippedCosmetics
	data, _ := json.Marshal(cur)
	os.WriteFile(saveFile, data, 0644)
}

func saveState(totalBells int, unlockedCosmetics []int, equippedCosmetics map[string]int) {
	cur := loadSave()
	cur.TotalBells = totalBells
	cur.UnlockedCosmetics = unlockedCosmetics
	cur.EquippedCosmetics = equippedCosmetics
	data, _ := json.Marshal(cur)
	os.WriteFile(saveFile, data, 0644)
}
