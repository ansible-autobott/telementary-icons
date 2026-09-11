package main

import "os"

// copyFile copies src to dst, following symlinks so an alias source becomes a
// real file at the destination.
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
