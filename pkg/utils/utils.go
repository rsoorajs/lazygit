package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/jesseduffield/gocui"
)

// GetProjectRoot returns the path to the root of the project. Only to be used
// in testing contexts, as with binaries it's unlikely this path will exist on
// the machine
func GetProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return strings.Split(dir, "lazygit")[0] + "lazygit"
}

func SortRange(x int, y int) (int, int) {
	if x < y {
		return x, y
	}
	return y, x
}

func AsJson(i interface{}) string {
	bytes, _ := json.MarshalIndent(i, "", "    ")
	return string(bytes)
}

// used to keep a number n between 0 and max, allowing for wraparounds
func ModuloWithWrap(n, max int) int {
	if max == 0 {
		return 0
	}

	if n >= max {
		return n % max
	} else if n < 0 {
		return max + n
	}
	return n
}

func FindStringSubmatch(str string, regexpStr string) (bool, []string) {
	re := regexp.MustCompile(regexpStr)
	match := re.FindStringSubmatch(str)
	return len(match) > 0, match
}

func MustConvertToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

// Safe will close tcell if a panic occurs so that we don't end up in a malformed
// terminal state
func Safe(f func()) {
	_ = SafeWithError(func() error { f(); return nil })
}

func SafeWithError(f func() error) error {
	panicking := true
	defer func() {
		if panicking && gocui.Screen != nil {
			gocui.Screen.Fini()
		}
	}()

	err := f()

	panicking = false

	return err
}

func StackTrace() string {
	buf := make([]byte, 10000)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s\n", buf[:n])
}

// returns the path of the file that calls the function.
// 'skip' is the number of stack frames to skip.
func FilePath(skip int) string {
	_, path, _, _ := runtime.Caller(skip)
	return path
}
