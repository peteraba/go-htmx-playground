package htmx

func IsHx(headers map[string][]string) bool {
	if len(headers) == 0 {
		return false
	}

	if v, ok := headers["Hx-Request"]; ok && len(v) > 0 {
		return v[0] == "true"
	}

	return false
}
