package project_evidence

import (
	"strings"
)

func validateRequired(fields map[string]string) []string {
	var missing []string
	for name, value := range fields {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	return missing
}