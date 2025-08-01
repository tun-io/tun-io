package headers

var disallowedHeaders = map[string]struct{}{
	"Connection": {},
}

func IsDisallowedHeader(header string) bool {
	_, exists := disallowedHeaders[header]
	return exists
}
