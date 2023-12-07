package htmx

const (
	HeaderHxRequest = "Hx-Request"
	HeaderHxTarget  = "Hx-Target"
)

func IsHx(headers map[string][]string) bool {
	if len(headers) == 0 {
		return false
	}

	v, ok := headers[HeaderHxRequest]
	if ok && len(v) > 0 {
		return v[0] == "true"
	}

	return false
}

func GetTarget(headers map[string][]string) string {
	v, ok := headers[HeaderHxTarget]
	if ok && len(v) > 0 {
		return v[0]
	}

	return ""
}
