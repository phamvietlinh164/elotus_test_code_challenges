package utils

func GetOffset(pageSize int, pageNumber int) int {
	if pageNumber <= 0 {
		return 0
	}
	return (pageNumber - 1) * pageSize
}
