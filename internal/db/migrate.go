package db

import (
	"os"
	"path/filepath"
	"strings"

	"xarr-proxy/internal/config"
	"xarr-proxy/internal/consts"
	"xarr-proxy/internal/utils"
)

func HotPatch(cfg *config.Config) error {
	// check if the checkpoint file exists
	var (
		version string
		err     error
	)

	checkpointPath := filepath.Join(cfg.ConfDir, consts.CHECKPOINT_FILENAME)
	if version, err = utils.ReadFile(checkpointPath); err == os.ErrNotExist {
		// create the checkpoint files
		v := consts.VERSION
		err := utils.CreateFile(checkpointPath, &v)
		if err != nil {
			return err
		}
	}
	files, err := utils.WalkDir(consts.SQL_FILE_DIR, func(filename string) bool {
		return strings.HasSuffix(filename, ".sql")
	})
	if err != nil {
		return err
	}
	for _, file := range files {
		// compare two strings
		if strings.Compare(version, file) < 0 {
			// execute the sql file
			if fileContent, err := utils.ReadFile(file); err != nil {
				return err
			} else {
				if err := db.Exec(fileContent).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Migrate(cfg *config.Config) {
	if err := HotPatch(cfg); err != nil {
		panic(err)
	}
}
