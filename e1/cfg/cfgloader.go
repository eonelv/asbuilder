package cfg

import (
	"os"
	"bufio"
	"strings"
	_ "github.com/mattn/go-sqlite3"
)
const (
	SERVER_PORT string = "SERVER_PORT"
	SERVER_HOME string = "SERVER_HOME"
	DB_NAME string = "DB_NAME"
	BUILD_CMD string = "BUILD_CMD"
	BUILD_CMD_SP string = "BUILD_CMD_SP"
)

var ServerCfg map[string]string
func LoadCfg() (bool, error) {
	userFile := "server.cfg"

	file,err := os.OpenFile(userFile, os.O_RDONLY, os.ModeAppend)

	if err != nil {
		return false, err
	}
	defer file.Close()

	ServerCfg = make(map[string]string)
	// using scanner to read config file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// expression match
		lineArray := strings.Split(line, "=")
		if len(lineArray) >= 2 {
			ServerCfg[lineArray[0]] = lineArray[1]
		}
	}
	return true, nil
}
