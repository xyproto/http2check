package vt100

// TODO: Find the terminal size on Windows
func TermSize() (uint, uint, error) {
	return 80, 25, nil
}
