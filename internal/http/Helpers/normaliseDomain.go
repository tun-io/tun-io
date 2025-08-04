package Helpers

func NormaliseUrl(url string) string {
	if len(url) < 8 {
		return url
	}

	if url[:7] == "http://" {
		return url[7:]
	}

	if url[:8] == "https://" {
		return url[8:]
	}

	return url
}

func AddHttpPrefix(url string) string {
	if len(url) < 7 {
		return "http://" + url
	}

	if url[:7] != "http://" && url[:8] != "https://" {
		return "http://" + url
	}

	return url
}
