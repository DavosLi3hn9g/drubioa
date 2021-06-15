package fun

func InSliceInt(value int, slice []int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func InSliceString(value string, slice []string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

//String切片去重
func SliceUniqueString(value *[]string) *[]string {
	result := make([]string, 0, len(*value))
	temp := map[string]struct{}{}
	for _, item := range *value {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return &result
}

//Int切片去重
func SliceUniqueInt(value *[]int) *[]int {
	result := make([]int, 0, len(*value))
	temp := map[int]struct{}{}
	for _, item := range *value {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return &result
}
