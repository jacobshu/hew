package shared

// returns items in a that aren't in b
func Difference(a, b []string) []string {
	set := make(map[string]struct{}, len(b))
	for _, val := range b {
		set[val] = struct{}{}
	}

	var diff []string
	for _, s := range a {
		if _, found := set[s]; !found {
			diff = append(diff, s)
		}
	}

	return diff
}
