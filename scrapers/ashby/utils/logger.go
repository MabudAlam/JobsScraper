package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"os"
	"strings"
	"time"
)

var LoggerInstance *Logger

func init() {
	LoggerInstance = &Logger{}
}

type Logger struct{}

func (l *Logger) Info(v ...any) {
	printLog("INFO", v...)
}

func (l *Logger) Error(v ...any) {
	printLog("ERROR", v...)
}

func (l *Logger) Warn(v ...any) {
	printLog("WARN", v...)
}

func (l *Logger) Debug(v ...any) {
	if os.Getenv("DEBUG") == "1" {
		printLog("DEBUG", v...)
	}
}

func printLog(level string, v ...any) {
	var sb strings.Builder
	for i, val := range v {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(toString(val))
	}
	println("[" + level + "] " + sb.String())
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return itoa(t)
	case int64:
		return itoa(int(t))
	default:
		return ""
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	if neg {
		s = "-" + s
	}
	return s
}

func ContentHash(fields ...string) string {
	h := sha256.New()
	for _, f := range fields {
		h.Write([]byte(f))
	}
	return hex.EncodeToString(h.Sum(nil))[:16]
}

func JitteredDelay(minMs, maxMs int) {
	delay := minMs + rand.Intn(maxMs-minMs)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
