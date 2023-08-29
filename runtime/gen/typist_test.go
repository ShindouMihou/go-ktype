package gen

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	err := Load("../../generators/kotlin/data_class.kt.typist")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSection_Replace(t *testing.T) {
	typist := Typists["../../generators/kotlin/data_class.kt.typist"]
	result := typist["package"].Replace(map[string]string{
		"{$PACKAGE}": "pro.qucy.ktype",
	})
	if result != "package pro.qucy.ktype" {
		t.Fatal("package does not match expected result, found: ", result)
	}
}

func TestCombine(t *testing.T) {
	typist := Typists["../../generators/kotlin/data_class.kt.typist"]
	fmt.Println(
		Combine(
			typist["package"].Replace(map[string]string{
				"{$PACKAGE}": "pro.qucy.ktype",
			}),
			DoubleBackedWhitespace,
			typist["imports"].String(),
			DoubleBackedWhitespace,
			typist["class_declaration"].Replace(map[string]string{
				"{$NAME}": "Server",
				"{$FIELDS}": Combine(
					typist["field_declaration"].Replace(map[string]string{
						"{$JSON_KEY}":   "id",
						"{$FIELD_NAME}": "id",
						"{$TYPE}":       "Long",
					}),
				),
			}),
		),
	)
}
