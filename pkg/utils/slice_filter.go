package utils

func FilterString(sliceToBeFiltered []string, test func(string) bool) (ret []string) {
	for _, value := range sliceToBeFiltered {
		if test(value) {
			ret = append(ret, value)
		}
	}
	return
}
