package basics

// Subtract returns all items in A (as) minus those present in B (bs)
func Subtract(as []string, bs []string) []string {
	setB := make(map[string]bool)
	for _, b := range bs {
		setB[b] = true
	}

	var diff []string
	for _, a := range as {
		_, existed := setB[a]
		if !existed {
			diff = append(diff, a)
		}
	}
	return diff
}
