package utility

func RemoveLastTwoItems[T any](slice []T) []T {
	if len(slice) < 2 {
		return []T{} // Return an empty slice if there are fewer than 2 items
	}
	return slice[:len(slice)-2]
}
