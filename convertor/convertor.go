package convertor

// int to base62 convertor
func Convertor(id int64) string {

	var chars string = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var res string = ""
	for id != 0 {
		res += string(chars[id%62])
		id = id / 62
	}

	for len(res) < 0 {
		res += "0"
	}
	return reverse(res)
}

func reverse(s string) string {
	str := []rune(s)

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}

	return string(str)
}
