package utils

import (
	"fmt"
	"strings"
)

var markdownV2Replacer = strings.NewReplacer(
	`\\`, `\\\\`,
	`_`, `\\_`,
	`*`, `\\*`,
	`[`, `\\[`,
	`]`, `\\]`,
	`(`, `\\(`,
	`)`, `\\)`,
	`~`, `\\~`,
	"`", "\\`",
	`>`, `\\>`,
	`#`, `\\#`,
	`+`, `\\+`,
	`-`, `\\-`,
	`=`, `\\=`,
	`|`, `\\|`,
	`{`, `\\{`,
	`}`, `\\}`,
	`.`, `\\.`,
	`!`, `\\!`,
)

func EscapeMarkdownV2(text string) string {
	if text == "" {
		return text
	}
	return markdownV2Replacer.Replace(text)
}

func FormatUintSlice(items []uint) string {
	if len(items) == 0 {
		return "-"
	}
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, fmt.Sprintf("%d", item))
	}
	return strings.Join(values, ", ")
}
