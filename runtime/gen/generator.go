package gen

import (
	"github.com/bytedance/sonic"
	"go-ktype/internal/utils"
	"go-ktype/runtime"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"os"
	"strings"
	"unicode"
)

type ImportMap = map[string]string
type Properties = map[string]string

type Transpiler struct {
	Typist     Typist
	ImportMap  ImportMap
	Properties Properties
	imports    []string
}

func NewTranspiler(typist Typist, importMap string, properties Properties) (*Transpiler, error) {
	var aliases map[string]string
	if importMap != "" {
		immf, err := os.Open(importMap)
		if err != nil {
			return nil, err
		}
		defer immf.Close()
		immb, err := io.ReadAll(immf)
		if err != nil {
			return nil, err
		}
		if err := sonic.Unmarshal(immb, &aliases); err != nil {
			return nil, err
		}
	}
	return &Transpiler{
		Typist:     typist,
		ImportMap:  aliases,
		Properties: properties,
		imports:    []string{},
	}, nil
}

func (transpiler *Transpiler) Transpile(runtime *runtime.Runtime) string {
	pkg := transpiler.Typist["package"].Replace(map[string]string{
		"{$PACKAGE}": transpiler.Properties["package"],
	})
	imports := strings.Builder{}
	classes := strings.Builder{}

	for _, class := range runtime.Classes {
		classes.WriteString(transpiler.encode(class))
		classes.WriteString(DoubleBackedWhitespace)
	}

	for _, imprt := range transpiler.imports {
		imports.WriteString(transpiler.Typist["import"].Replace(map[string]string{
			"{$PACKAGE}": imprt,
		}))
		imports.WriteString("\n")
	}

	return Combine(
		pkg,
		DoubleBackedWhitespace,
		imports.String(),
		transpiler.Typist["imports"].String(),
		DoubleBackedWhitespace,
		classes.String(),
	)
}

func (transpiler *Transpiler) encode(class runtime.Class) string {
	fields := strings.Builder{}
	fCount := 0
	for field, t := range class.Fields {
		if alias, ok := transpiler.ImportMap[t]; ok {
			if !utils.AnyMatchStringCaseInsensitive(transpiler.imports, alias) {
				transpiler.imports = append(transpiler.imports, alias)
			}
		}

		properties := class.Properties[field]
		nullable := false

		if _, ok := properties["type_gen"]; !ok {
			t = cases.Title(language.English).String(t)
		}

		if utils.HasPrefixStr(t, "Nullable[") {
			t, _ = strings.CutPrefix(t, "Nullable[")
			t, _ = strings.CutSuffix(t, "]")
			nullable = true
		}

		if utils.HasPrefixStr(t, "Slice[") {
			t, _ = strings.CutPrefix(t, "Slice[")
			t, _ = strings.CutSuffix(t, "]")
			t = "List<" + t + ">"
		}
		key := field
		if k, ok := properties["json"]; ok {
			key = k
		}
		if nullable {
			t += "?"
		}

		if unicode.IsUpper(rune(field[0])) {
			field = string(unicode.ToLower(rune(field[0]))) + field[1:]
		}

		fields.WriteString(transpiler.Typist["field_declaration"].Replace(map[string]string{
			"{$JSON_KEY}":   key,
			"{$FIELD_NAME}": field,
			"{$TYPE}":       t,
		}))

		if fCount < (len(class.Fields) - 1) {
			fields.WriteString(",\n")
		}

		fCount++
	}
	return Combine(
		transpiler.Typist["class_declaration"].Replace(map[string]string{
			"{$NAME}":   class.Name,
			"{$FIELDS}": fields.String(),
		}),
	)
}
