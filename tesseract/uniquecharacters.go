package tesseract

func (me *Tesseract) uniqueCharacters(str string) string {
	charSet := make(map[rune]bool)

	for _, char := range str {
		charSet[char] = true
	}

	keys := make([]rune, 0, len(charSet))
	res := ""
	for k := range charSet {
		keys = append(keys, k)
		if (k >= 65 && k <= 90) || (k >= 97 && k <= 122) || (k >= 48 && k <= 57) || (k == 252) || (k == 246) || (k == 228) {
			res += string(k)
		}
	}

	return res
}
