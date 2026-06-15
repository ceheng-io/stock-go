package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// FetchJSVarsOptions configures FetchJSVars.
type FetchJSVarsOptions struct {
	Timeout time.Duration
	Headers map[string]string
}

// ExtractJSVar extracts and decodes one JSON-compatible JavaScript variable declaration.
func ExtractJSVar(text string, name string, target any) error {
	pattern := regexp.MustCompile(`(?:^|[^\w$])(?:var|let|const)\s+` + regexp.QuoteMeta(name) + `\s*=\s*`)
	match := pattern.FindStringIndex(text)
	if match == nil {
		return fmt.Errorf("js variable %q not found", name)
	}
	start := match[1]
	end := jsVarValueEnd(text, start)
	literal := text[start:end]
	if err := json.Unmarshal([]byte(literal), target); err != nil {
		return fmt.Errorf("invalid js variable %q: %w", name, err)
	}
	return nil
}

// ParseJSVars extracts JSON-compatible JavaScript variable declarations.
func ParseJSVars(text string, names ...string) map[string]any {
	values := make(map[string]any)
	for _, name := range names {
		var value any
		if err := ExtractJSVar(text, name, &value); err == nil {
			values[name] = value
		}
	}
	return values
}

// FetchJSVars fetches a JavaScript variable declaration document and extracts variables.
func FetchJSVars(ctx context.Context, client *http.Client, requestURL string, names []string, options FetchJSVarsOptions) (map[string]any, error) {
	if client == nil {
		client = http.DefaultClient
	}
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer drainAndClose(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, HTTPStatusError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			URL:        requestURL,
		}
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return ParseJSVars(string(body), names...), nil
}

func jsVarValueEnd(text string, start int) int {
	depth := 0
	var quote byte
	escaped := false
	for index := start; index < len(text); index++ {
		ch := text[index]
		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			if ch == '\\' {
				escaped = true
				continue
			}
			if ch == quote {
				quote = 0
			}
			continue
		}
		switch ch {
		case '"', '\'':
			quote = ch
		case '[', '{', '(':
			depth++
		case ']', '}', ')':
			if depth > 0 {
				depth--
			}
		case ';':
			if depth == 0 {
				return index
			}
		}
	}
	return len(text)
}
