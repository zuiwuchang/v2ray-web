package utils

import "syscall"

// Quit send quit signal
func Quit() {
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)
}
