package tests

import (
	"fmt"
	"strings"
	"testing"

	lst "github.com/binchencoder/letsgo/strings"
)

func TestReplaceAll(t *testing.T) {
	oldStr := "vexillary-service=vxserver:4100"
	fmt.Printf("oldString: %s\n, newString: %s\n", oldStr, strings.ReplaceAll(oldStr, "=", "/"))
}

func TestCsvString(t *testing.T) {
	str := "192.168.38.6:1900,192.168.39.6:1900"
	fmt.Printf("str: %s\n, csvString: %s\n", str, lst.CsvToSlice(str))
}
