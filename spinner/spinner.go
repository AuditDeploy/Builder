package spinner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/theckman/yacspin"
	"go.uber.org/zap"
)

var BuilderLog = zap.S()

var cfg = yacspin.Config{
	Frequency:     100 * time.Millisecond,
	CharSet:       []string{"⣧", "⣇", "⡇", "⠇", "⠃", "⠁", "⠃", "⠇", "⡇", "⣇", "⣧", "⣧"},
	StopCharacter: "",
}
var Spinner, err = yacspin.New(cfg)
var Caller string

func LogMessage(msg string, level string) {
	args := os.Args[1:]
	_, file, no, ok := runtime.Caller(1)
	if ok {
		Caller = filepath.Base(file) + ":" + fmt.Sprint(no)
	}

	Spinner.Stop()

	// Check if debug flag is given
	for i := 0; i < len(args); i++ {
		if args[i] == "-d" || args[i] == "--debug" {
			// Print log message at correct level
			switch level {
			case "info":
				fmt.Println(time.Now().Local().String() + "   INFO   " + Caller + ":   " + msg)
				break
			case "warn":
				fmt.Println(time.Now().Local().String() + "   WARN   " + Caller + ":   " + msg)
				break
			case "error":
				fmt.Println(time.Now().Local().String() + "   ERROR   " + Caller + ":   " + msg)
				break
			default: // Fatal
				fmt.Println(time.Now().Local().String() + "   FATAL   " + Caller + ":   " + msg)
				BuilderLog.Fatal()
			}
		}
	}

	Spinner.Start()
}
