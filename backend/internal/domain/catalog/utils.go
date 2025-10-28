package catalog

import (
	"regexp"
	"strings"
)

var extensionRegex = regexp.MustCompile(`\.(csv|xlsx|xls|tsv|txt)$`)

// NormalizeSourceName normalizes a source name for case-insensitive matching
// Examples:
//   "Cannabis_Inventory.csv" -> "cannabis_inventory"
//   "SALES DATA.CSV" -> "sales data"
//   "Products.xlsx" -> "products"
func NormalizeSourceName(name string) string {
	// Remove file extension
	normalized := extensionRegex.ReplaceAllString(name, "")

	// Trim whitespace and convert to lowercase
	normalized = strings.ToLower(strings.TrimSpace(normalized))

	return normalized
}
