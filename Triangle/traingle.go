package Triangle

func traingle(s1, s2, s3 int) string {
	if s1+s2 >= s3 || s1+s3 >= s2 || s2+s3 >= s1 {
		if s1 == s2 && s2 == s3 {
			return "ISO"
		} else if s1 == s2 || s2 == s3 || s1 == s3 {
			return "Eque"
		} else if (s1 != s2 && s2 != s3) && (s1 != s3) {
			return "scalar"
		}

	}

	return "not valid"
}
