package utils

// Contains checks if a slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsFlag checks if a slice contains a specific flag (alias for Contains for clarity)
func ContainsFlag(args []string, flag string) bool {
	return Contains(args, flag)
}
