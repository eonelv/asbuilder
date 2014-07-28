package log

import (
	"os"
	"log"
	"fmt"
)

type ELogger struct {
	Log *log.Logger
	logFile *os.File
}

var LoggerSys *ELogger

func CreateLogger(name string) (*ELogger, error) {
	logfile, err:=os.OpenFile(name,os.O_APPEND|os.O_CREATE, os.FileMode(0666))
	if err != nil{
		return nil, err
	}
	logger:=log.New(logfile,"",log.Ldate|log.Ltime|log.Llongfile)

	elogger := &ELogger{}
	elogger.Log = logger

	return elogger, nil
}

func (this *ELogger) Close() {
	this.logFile.Close()
}

func (l *ELogger) Printf(format string, v ...interface{}) {
	l.Log.Output(2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *ELogger) Print(v ...interface{}) { l.Log.Output(2, fmt.Sprint(v...)) }

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *ELogger) Println(v ...interface{}) { l.Log.Output(2, fmt.Sprintln(v...)) }

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *ELogger) Fatal(v ...interface{}) {
	l.Log.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *ELogger) Fatalf(format string, v ...interface{}) {
	l.Log.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *ELogger) Fatalln(v ...interface{}) {
	l.Log.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}
