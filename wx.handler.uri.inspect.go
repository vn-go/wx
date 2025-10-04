package wx

import (
	"regexp"
	"strings"
)

type inspectUriInfo struct {
	// mainUri is the base path used for mapping handlers.
	// 1. For "files/{*FilePath}", mainUri is 'files'.
	// 2. For "file?path={path}&name={name}", mainUri is 'file'.
	// 3. For "file/{path}/{name}", mainUri is 'file'.
	mainUri string
	// swaggerUri is the standardized URI template used for documentation (Swagger/OpenAPI).
	// 1. For "files/*path", swaggerUri is "files/{path}".
	// 2. For "/media/files3/{FilePath}/{filename}?test={code}", swaggerUri is the same.
	swaggerUri string
	// uriRegex is the regular expression used by the routing engine to match incoming requests
	// and extract path parameters.
	// 1. For "media/files/{*FilePath}" -> uriRegex = "media/files/(.*)"
	// 2. For "media/files2/{FilePath}/{filename}" -> uriRegex = "^media/files2/([^/]+)/([^/]+)$"
	// 3. For "media/files3/{FilePath}/{filename}?test={code}" -> uriRegex = "^media/files3/([^/]+)/([^/]+)$" (captures up to the first '?')
	uriRegex string
	// isRegex is set to true if the URI contains path parameters that require a regex match
	// (i.e., contains "{" or "*").
	isRegex     bool
	uriParams   []string
	queryParams []string
}

// placeholderRegex matches the two types of path placeholders: {param} and {*param}

// Hỗ trợ {param}, {*param}, *param
var placeholderRegex = regexp.MustCompile(`\{[*]?[^}]+\}|\*[^/]+`)

func inspectUri(fullUrl string) inspectUriInfo {
	uri := fullUrl

	info := inspectUriInfo{
		swaggerUri:  uri,
		uriParams:   []string{},
		queryParams: []string{},
	}
	info.mainUri = uri
	if strings.Contains(fullUrl, "?") {
		uri = strings.Split(fullUrl, "?")[0]
		items := strings.Split(strings.Split(fullUrl, "?")[1], ",")
		for _, x := range items {
			if strings.Contains(x, "=") {
				info.queryParams = append(info.queryParams, strings.Split(x, "=")[0])
			} else {
				info.queryParams = append(info.queryParams, x)
			}
		}
	}
	// --- 1. determine main uri ---
	pathPart := uri
	if qIndex := strings.Index(uri, "?"); qIndex != -1 {
		pathPart = uri[:qIndex]
	}

	// --- 2. Nếu không có { hoặc * thì không cần regex ---
	if !strings.Contains(uri, "{") && !strings.Contains(uri, "*") {
		info.isRegex = false
		return info
	} else {
		info.mainUri = ""
		segments := strings.Split(strings.Trim(pathPart, "/"), "/")
		for _, x := range segments {
			if strings.Contains(x, "{") {
				break
			} else {
				info.mainUri += x + "/"
			}
		}

		// if len(segments) > 0 && segments[0] != "" {
		// 	info.mainUri = segments[0]
		// } else {
		// 	info.mainUri = strings.Trim(uri, "/")
		// }
		// info.mainUri = strings.Split(uri, "{")[0]
	}

	// --- 3. Xây swaggerUri và uriRegex ---
	swaggerBuilder := strings.Builder{}
	regexBuilder := strings.Builder{}
	regexBuilder.WriteString("^")

	lastIndex := 0
	matches := placeholderRegex.FindAllStringIndex(pathPart, -1)

	for _, match := range matches {
		start, end := match[0], match[1]
		literal := pathPart[lastIndex:start]
		regexBuilder.WriteString(regexp.QuoteMeta(literal))
		swaggerBuilder.WriteString(literal)

		placeholder := pathPart[start:end]

		// --- Xác định tên tham số thật ---
		name := extractParamName(placeholder)
		info.uriParams = append(info.uriParams, name)

		if strings.Contains(placeholder, "*") {
			// Catch-all: {*FilePath} hoặc *FilePath
			swaggerBuilder.WriteString("{" + name + "}")
			regexBuilder.WriteString("(.*)")
		} else {
			// Normal placeholder: {filename}
			swaggerBuilder.WriteString("{" + name + "}")
			regexBuilder.WriteString("([^/]+)")
		}

		lastIndex = end
	}

	// --- 4. Thêm phần literal còn lại ---
	if lastIndex < len(pathPart) {
		literal := pathPart[lastIndex:]
		regexBuilder.WriteString(regexp.QuoteMeta(literal))
		swaggerBuilder.WriteString(literal)
	}

	regexBuilder.WriteString("$")

	// --- 5. Gán kết quả ---
	info.uriRegex = regexBuilder.String()
	info.swaggerUri = swaggerBuilder.String()

	// Nếu có query string, nối thêm vào swaggerUri (Swagger có thể mô tả query params riêng)
	// if qIndex := strings.Index(uri, "?"); qIndex != -1 {
	// 	info.swaggerUri += uri[qIndex:]
	// 	items := strings.Split(uri[qIndex+1:], ",")
	// 	for _, x := range items {
	// 		if strings.Contains(x, "=") {
	// 			items = append(items, strings.Split(x, "=")[0])
	// 		} else {
	// 			items = append(items, x)
	// 		}
	// 	}
	// }

	info.isRegex = true
	return info
}
func extractParamName(s string) string {
	s = strings.Trim(s, "{}")
	s = strings.TrimPrefix(s, "*")
	s = strings.TrimSpace(s)
	if s == "" {
		return "param"
	}
	return s
}
