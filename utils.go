package sqlsqbdr

func MultipleStringToMap(ss []string) map[string]string {
	sm := make(map[string]string)
	for _, s := range ss {
		sm[s] = s
	}
	return sm
}
