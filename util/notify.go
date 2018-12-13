package util

import (
	//"runtime/debug"
	//"fmt"
	"testing"
)

func TFatal(t *testing.T, msg string) {
	//debug.PrintStack()
	DingNotify(t.Name() + ": " + msg)
	t.Fatal(msg)
}
