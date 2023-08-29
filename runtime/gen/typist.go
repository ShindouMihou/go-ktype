package gen

import (
	"bufio"
	"go-ktype/internal/utils"
	"io"
	"os"
	"strings"
)

type Typist map[string]Section
type Section string

var Typists = make(map[string]Typist)
var DoubleBackedWhitespace = "\n\n"

func Load(files ...string) error {
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		typist := NewTypist(f)
		Typists[file] = typist
		_ = f.Close()
	}
	return nil
}

func NewTypist(rd io.Reader) Typist {
	typist := make(Typist)

	buf := bufio.NewScanner(rd)
	currentSection := "default"
	builder := strings.Builder{}

	for buf.Scan() {
		text := buf.Text()
		if utils.HasPrefixStr(text, "--typist:section ") {
			if currentSection != "default" {
				typist[currentSection] = Section(builder.String())
				builder.Reset()
			}
			currentSection = strings.SplitN(text, " ", 2)[1]
			continue
		}
		if builder.Len() != 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(text)
	}
	typist[currentSection] = Section(builder.String())
	builder.Reset()

	return typist
}

func (section Section) String() string {
	return string(section)
}

func (section Section) Replace(modifiers map[string]string) string {
	cpy := string(section)
	for key, value := range modifiers {
		cpy = strings.ReplaceAll(cpy, key, value)
	}
	return cpy
}

func Combine(sections ...string) string {
	combined := ""
	for _, section := range sections {
		combined += section
	}
	return combined
}
