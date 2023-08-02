package utils

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func GetGoroutineID() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// Parse the 4708504 out of "goroutine 4708504 [running]:"
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return fmt.Sprintf("%d", id)
}
