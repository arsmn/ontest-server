package daemon

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/arsmn/ontest/module/httplib"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func startupMessage(addr string, tls bool, handlersCount, templatesCount uint32, cfgFile string) {
	var logo string
	logo += "\n%s"
	logo += " ┌─────────────────────────────────────────────────────┐\n"
	logo += " │ %s   │\n"
	logo += " │ %s   │\n"
	logo += " │                                                     │\n"
	logo += " │ Handlers %s  Templates %s │\n"
	logo += " │ Go ......%s  Threads ..%s │\n"
	logo += " │ OS ......%s  PID ......%s │\n"
	logo += " │ Config file %s \n"
	logo += " └─────────────────────────────────────────────────────┘"
	logo += "%s\n\n"

	const (
		cBlack   = "\u001b[90m"
		cRed     = "\u001b[91m"
		cCyan    = "\u001b[96m"
		cGreen   = "\u001b[92m"
		cYellow  = "\u001b[93m"
		cBlue    = "\u001b[94m"
		cMagenta = "\u001b[95m"
		cWhite   = "\u001b[97m"
		cReset   = "\u001b[0m"
	)

	value := func(s string, width int) string {
		pad := width - len(s)
		str := ""
		for i := 0; i < pad; i++ {
			str += "."
		}
		str += fmt.Sprintf(" %s%s%s", cCyan, s, cBlack)
		return str
	}

	center := func(s string, width int) string {
		pad := strconv.Itoa((width - len(s)) / 2)
		str := fmt.Sprintf("%"+pad+"s", " ")
		str += s
		str += fmt.Sprintf("%"+pad+"s", " ")
		if len(str) < width {
			str += " "
		}
		return str
	}

	centerValue := func(s string, width int) string {
		pad := strconv.Itoa((width - len(s)) / 2)
		str := fmt.Sprintf("%"+pad+"s", " ")
		str += fmt.Sprintf("%s%s%s", cCyan, s, cBlack)
		str += fmt.Sprintf("%"+pad+"s", " ")
		if len(str)-10 < width {
			str += " "
		}
		return str
	}

	host, port := httplib.ParseAddr(addr)
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	addr = "http://" + host + ":" + port
	if tls {
		addr = "https://" + host + ":" + port
	}

	out := colorable.NewColorableStdout()
	if os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		out = colorable.NewNonColorable(os.Stdout)
	}
	fmt.Fprintf(out, logo,
		cBlack,
		centerValue(" On-Test", 49),
		center(addr, 49),
		value(strconv.Itoa(int(handlersCount)), 14), value(strconv.Itoa(int(templatesCount)), 14),
		value(runtime.Version(), 14), value(strconv.Itoa(runtime.NumCPU()), 14),
		value(runtime.GOOS, 14), value(strconv.Itoa(os.Getpid()), 14),
		value(cfgFile, 14),
		cReset,
	)
}
