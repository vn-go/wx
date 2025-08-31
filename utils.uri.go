package wx

import (
	"regexp"
	"strings"

	"github.com/vn-go/wx/internal"
)

type uriHelperType struct {
	SpecialCharForRegex string
}

/*
handlerInfoExtractUriParams extracts all substrings enclosed in curly braces '{}'
from the given URI string, along with their positions (based on segments split by '/').

For example:

Given URI: "abc/{word 1}/abc/dbc/{u2}/{u3}"

The function will return a slice of uriParaim structs:
[

	{Position: 1, Name: "word 1"},
	{Position: 4, Name: "u2"},
	{Position: 5, Name: "u3"},

]

Where:
- Position is the zero-based index of the segment in the URI path split by '/'.
- Name is the trimmed string inside the braces '{}'.

@return []uriParaim - a slice containing extracted parameters with their position and name.
*/
func (h *uriHelperType) ExtractUriParams(uri string) []uriParam {
	key := uri + "/helperType/ExtractUriParams"
	ret, _ := internal.OnceCall(key, func() (*[]uriParam, error) {

		params := []uriParam{}
		segments := h.SplitUriSegments(strings.Split(uri, "?")[0])

		for i, segment := range segments {
			// Check if segment contains a URI parameter enclosed in {}
			name := h.ExtractNameInBraces(segment)
			if name != "" {
				isSlug := false

				if name[0] == '*' {
					name = name[1:]
					isSlug = true
				}
				params = append(params, uriParam{
					Position: i,
					Name:     name,
					IsSlug:   isSlug,
				})
			}
		}

		return &params, nil
	})
	return *ret
}

// splitUriSegments splits the URI string by '/', ignoring empty segments.
func (h *uriHelperType) SplitUriSegments(uri string) []string {
	var segments []string
	start := 0

	for i := 0; i < len(uri); i++ {
		if uri[i] == '/' {
			if start < i {
				segments = append(segments, uri[start:i])
			}
			start = i + 1
		}
	}
	// append the last segment if any
	if start < len(uri) {
		segments = append(segments, uri[start:])
	}
	return segments
}

// extractNameInBraces extracts the trimmed content inside the first pair of braces '{}' in the segment.
// Returns empty string if no braces found.
func (h *uriHelperType) ExtractNameInBraces(segment string) string {
	start := -1
	end := -1
	for i, ch := range segment {
		if ch == '{' && start == -1 {
			start = i
		} else if ch == '}' && start != -1 {
			end = i
			break
		}
	}
	if start != -1 && end != -1 && end > start+1 {
		name := segment[start+1 : end]
		return h.TrimSpaces(name)
	}
	return ""
}

// trimSpaces trims leading and trailing spaces from a string.
func (h *uriHelperType) TrimSpaces(s string) string {
	start, end := 0, len(s)-1
	for start <= end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end >= start && (s[end] == ' ' || s[end] == '\t') {
		end--
	}
	if start > end {
		return ""
	}
	return s[start : end+1]
}
func (h *uriHelperType) calculateUrlWithQuery(ret *handlerInfo) {
	ret.queryParams = []queryParam{}

	uri := strings.TrimSuffix(strings.Split(ret.uri, "?")[0], "/")
	ret.uriQuery = strings.Split(ret.uri, "?")[1]
	ret.uri = uri

	//ret.uriHandler = strings.TrimSuffix(strings.Split(uri, "?")[0], "/")
	items := strings.Split(ret.uriQuery, "&")
	for _, x := range items {
		fieldName := strings.Split(x, "=")[1]
		fieldName = strings.TrimPrefix(fieldName, "{")
		fieldName = strings.TrimSuffix(fieldName, "}")
		field, ok := ret.typeOfArgIsHttpContextElem.FieldByNameFunc(func(s string) bool {
			return strings.EqualFold(s, fieldName)
		})
		if !ok {
			continue
		}
		ret.queryParams = append(ret.queryParams, queryParam{
			Name:       fieldName,
			FieldIndex: field.Index,
		})
	}

}
func (h *uriHelperType) convertUrlToRegex(urlPattern string) string {
	// Bước 1: Thay thế các wildcard catch-all (*...) bằng .*
	//regexPattern := strings.ReplaceAll(urlPattern, "*", ".*")

	// Bước 2: Xử lý các tham số đường dẫn thông thường {name}
	// Ví dụ: {id} sẽ được chuyển thành ([^/]+) để khớp với bất kỳ ký tự nào ngoại trừ "/"
	// re := regexp.MustCompile(`{[^}]+}`)
	re := regexp.MustCompile(`\{[*][^{}]+\}`)
	reParam := regexp.MustCompile(`\{[^{}*]+\}`)
	// Thay thế tất cả các khớp với biểu thức (.*)
	// Lưu ý: Chúng ta dùng (.*) để bắt toàn bộ nội dung, bao gồm cả dấu gạch chéo
	regexPattern := re.ReplaceAllString(urlPattern, "(.*)")
	regexPattern = reParam.ReplaceAllString(regexPattern, "([^/]+)")

	return regexPattern
}
func (h *uriHelperType) calculateUrl(ret *handlerInfo) {
	if len(ret.uriParams) > 0 {

		if !strings.Contains(ret.uri, "{*") {
			ret.regexUri = h.TemplateToRegex(ret.uri)
			ret.uriHandler = strings.Split(ret.uri, "{")[0]
		} else {
			ret.regexUri = h.convertUrlToRegex(ret.uri)
			ret.uriHandler = strings.Split(ret.uri, "{")[0]
		}

		ret.isRegexHandler = true
		ret.regexUriFind = *regexp.MustCompile(strings.ReplaceAll(strings.TrimPrefix(ret.regexUri, "^"), "/", "\\/"))

	} else {
		ret.regexUri = h.EscapeSpecialCharsForRegex(ret.uri)
		if ret.isRegexHandler {
			ret.uriHandler = ret.uri + "/"
		} else {
			ret.uriHandler = ret.uri
		}
	}
}

// templateToRegex chuyển URI template thành regex pattern string
// lấy các giá trị trong {}
func (h *uriHelperType) TemplateToRegex(template string) string {
	key := template + "/helperType/TemplateToRegex"
	ret, _ := internal.OnceCall(key, func() (*string, error) {
		segments := strings.Split(template, "/")
		regexParts := []string{}
		paramCount := 0
		var escapeRegex = h.EscapeSpecialCharsForRegex
		for _, seg := range segments {
			if seg == "" {
				continue
			}

			var sb strings.Builder
			i := 0
			for i < len(seg) {
				start := strings.Index(seg[i:], "{")
				if start == -1 {
					// No more '{', escape remainder
					sb.WriteString(escapeRegex(seg[i:]))
					break
				}

				start += i
				end := strings.Index(seg[start:], "}")
				if end == -1 {
					// No closing brace, treat literally
					sb.WriteString(escapeRegex(seg[i:]))
					break
				}
				end += start

				// Escape static part before {
				if start > i {
					sb.WriteString(escapeRegex(seg[i:start]))
				}

				// Add capture group for parameter inside {}
				sb.WriteString(`([^/]+)`)
				paramCount++

				// Move index past "}"
				i = end + 1
			}

			regexParts = append(regexParts, sb.String())
		}

		// Join parts with '/'
		regexPattern := "^" + strings.Join(regexParts, "/") + "$"
		return &regexPattern, nil
	})
	return *ret
}
func (h *uriHelperType) EscapeSpecialCharsForRegex(s string) string {
	ret := ""
	for _, c := range s {
		if strings.Contains(h.SpecialCharForRegex, string(c)) {
			ret += "\\"
		}
		ret += string(c)
	}
	return ret
}
