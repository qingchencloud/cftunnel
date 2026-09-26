// Package presets defines the common one-click scenarios shown by the
// desktop client and exposed by the terminal program.
package presets

import "strings"

// Preset is a named local service scenario.
type Preset struct {
	ID    string
	Label string
	Hint  string
	Port  string
}

var all = []Preset{
	{ID: "frontend", Label: "前端预览", Hint: "Vite / Webpack", Port: "5173"},
	{ID: "webapp", Label: "Web 应用", Hint: "Next.js / React", Port: "3000"},
	{ID: "webhook", Label: "Webhook", Hint: "回调与接口调试", Port: "8080"},
	{ID: "homeassistant", Label: "Home Assistant", Hint: "本地控制台", Port: "8123"},
}

// All returns a copy so callers cannot mutate the registry.
func All() []Preset { return append([]Preset(nil), all...) }

// Find resolves a preset by id, case-insensitively.
func Find(id string) (Preset, bool) {
	for _, item := range all {
		if strings.EqualFold(item.ID, id) {
			return item, true
		}
	}
	return Preset{}, false
}
