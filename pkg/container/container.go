package container

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Container struct {
	ID      string
	Name    string
	Image   string
	Command string
	Status  string
}

func (c *Container) PrepareBundle() error {
	bundleDir := filepath.Join("/tmp/containerhub", c.ID)
	if err := os.MkdirAll(bundleDir, 0755); err != nil {
		return err
	}

	config := map[string]interface{}{
		"ociVersion": "1.0.2",
		"process": map[string]interface{}{
			"terminal": false,
			"user":     map[string]int{"uid": 0, "gid": 0},
			"args":     []string{"sh", "-c", c.Command},
			"env":      []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		},
		"root": map[string]interface{}{
			"path": c.Image,
		},
		"hostname": c.Name,
		"mounts":   []interface{}{},
	}

	configFile, err := os.Create(filepath.Join(bundleDir, "config.json"))
	if err != nil {
		return err
	}
	defer configFile.Close()

	return json.NewEncoder(configFile).Encode(config)
}
