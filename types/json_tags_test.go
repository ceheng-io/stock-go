package types

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"reflect"
	"strings"
	"testing"
)

func TestExportedStructFieldsHaveLowerCamelJSONTags(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", func(info fs.FileInfo) bool {
		name := info.Name()
		return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
	}, 0)
	if err != nil {
		t.Fatalf("parse types package: %v", err)
	}

	pkg := pkgs["types"]
	if pkg == nil {
		t.Fatal("types package not found")
	}

	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}
			for _, spec := range gen.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				if !typeSpec.Name.IsExported() {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				for _, field := range structType.Fields.List {
					if len(field.Names) == 0 {
						continue
					}
					for _, name := range field.Names {
						if !name.IsExported() {
							continue
						}
						if field.Tag == nil {
							t.Fatalf("%s.%s missing json tag", typeSpec.Name.Name, name.Name)
						}
						tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get("json")
						if tag == "" {
							t.Fatalf("%s.%s missing json tag", typeSpec.Name.Name, name.Name)
						}
						if tag != expectedJSONName(name.Name) {
							t.Fatalf("%s.%s json tag = %q, want %q", typeSpec.Name.Name, name.Name, tag, expectedJSONName(name.Name))
						}
					}
				}
			}
		}
	}
}

func expectedJSONName(name string) string {
	parts := splitIdentifier(name)
	if len(parts) == 0 {
		return name
	}
	var builder strings.Builder
	builder.WriteString(strings.ToLower(parts[0]))
	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		if part == "YoY" {
			builder.WriteString(part)
			continue
		}
		lower := strings.ToLower(part)
		builder.WriteString(strings.ToUpper(lower[:1]))
		builder.WriteString(lower[1:])
	}
	return builder.String()
}

func splitIdentifier(name string) []string {
	initialisms := []string{
		"CFFEX", "COMEX", "CSRC",
		"HTML", "HTTP", "JSON",
		"ETF", "API", "BPS", "EPS", "MACD", "NAV", "OBV", "ROC", "ROA", "ROE", "RSI", "SAR", "SDK", "THS", "URL",
		"CN", "EN", "HK", "ID", "PB", "PE", "SH", "SZ", "TZ", "US", "YoY",
		"A", "H", "ZT",
	}

	var parts []string
	for len(name) > 0 {
		if part, ok := consumeInitialism(name, initialisms); ok {
			parts = append(parts, part)
			name = name[len(part):]
			continue
		}

		end := 1
		for end < len(name) && !startsWithInitialism(name[end:], initialisms) {
			prev := name[end-1]
			curr := name[end]
			if isUpper(curr) && !isUpper(prev) && !isDigit(prev) {
				break
			}
			if isDigit(curr) != isDigit(prev) {
				break
			}
			end++
		}
		parts = append(parts, name[:end])
		name = name[end:]
	}
	return parts
}

func consumeInitialism(name string, initialisms []string) (string, bool) {
	for _, initialism := range initialisms {
		if !strings.HasPrefix(name, initialism) {
			continue
		}
		if len(initialism) == 1 && (len(name) == 1 || !isUpper(name[1])) {
			continue
		}
		if len(initialism) == 1 || initialism == "ZT" || initialism == "YoY" || len(initialism) > 1 {
			return initialism, true
		}
	}
	return "", false
}

func startsWithInitialism(name string, initialisms []string) bool {
	_, ok := consumeInitialism(name, initialisms)
	return ok
}

func isUpper(value byte) bool {
	return value >= 'A' && value <= 'Z'
}

func isDigit(value byte) bool {
	return value >= '0' && value <= '9'
}
