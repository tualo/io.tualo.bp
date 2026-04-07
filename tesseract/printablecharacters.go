package tesseract

func (me *Tesseract) printableCharacters(str string) string {
	res := ""
	for _, char := range str {
		if (char >= 65 && char <= 90) || (char >= 97 && char <= 122) || (char >= 48 && char <= 57) || (char == 252) || (char == 246) || (char == 228) {
			res += string(char)
		}
	}
	return res
}
