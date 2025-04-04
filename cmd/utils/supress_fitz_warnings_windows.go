//go:build windows
// +build windows

package utils

func SuppressFitzWarnings() {
	// No-op on Windows.
	// stderr redirection is not needed or would require WinAPI calls.
}
