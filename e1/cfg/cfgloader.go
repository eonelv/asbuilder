package cfg

import (
	"os"
	"bufio"
	"fmt"
	"io"
	"strings"
	_ "github.com/mattn/go-sqlite3"
	"e1/log"
)
const (
	SERVER_PORT string = "SERVER_PORT"
	SERVER_HOME string = "SERVER_HOME"
	DB_NAME string = "DB_NAME"
	BUILD_CMD string = "BUILD_CMD"
	BUILD_CMD_SP string = "BUILD_CMD_SP"
)

var ServerCfg map[string]string
func LoadCfg() bool{
	userFile := "server.cfg"

	file,err := os.OpenFile(userFile, os.O_RDONLY, os.ModeAppend)

	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	ServerCfg = make(map[string]string)
	for {
		lineDatas, _, err := reader.ReadLine()
		if err != nil && err != io.EOF {
			return false
		}
		line := string(lineDatas)
		lineArray := strings.Split(line, "=")
		if len(lineArray) >= 2 {
			ServerCfg[lineArray[0]] = lineArray[1]
		}
		if err != nil && err == io.EOF {
			break
		}
	}
	log.LoggerSys.Print("Server config load success.")
	return true
}
