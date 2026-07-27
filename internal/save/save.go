// Package save persists the high score (max wave reached) to save.json next to
// the executable (GDD §8).
package save

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type data struct {
	HighWave int `json:"highWave"`
}

func path() string {
	if exe, err := os.Executable(); err == nil {
		return filepath.Join(filepath.Dir(exe), "save.json")
	}
	return "save.json"
}

// Load reads the stored high score, or 0 if there is none.
func Load() int {
	b, err := os.ReadFile(path())
	if err != nil {
		return 0
	}
	var d data
	if json.Unmarshal(b, &d) != nil {
		return 0
	}
	return d.HighWave
}

// Save writes the high score, keeping the larger of the stored and given value.
func Save(highWave int) {
	if highWave <= Load() {
		return
	}
	if b, err := json.Marshal(data{HighWave: highWave}); err == nil {
		_ = os.WriteFile(path(), b, 0o644)
	}
}
