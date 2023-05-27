package db

import (
	"fmt"
	"strings"

	"xarr-proxy/internal/utils"

	"github.com/rs/zerolog/log"
)

const (
	SQL_FILE_DIR = "resources/sql"
)

func HotPatch() error {
	files, err := utils.WalkDir(SQL_FILE_DIR, func(filename string) bool {
		return strings.HasSuffix(filename, ".sql")
	})
	if err != nil {
		return err
	}
	log.Info().Msgf("Found %d files in %s", len(files), SQL_FILE_DIR)
	fmt.Println(files)

	return nil
}

func Migrate() {
	HotPatch()
}
