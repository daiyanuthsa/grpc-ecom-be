package utils

import (
	"fmt"
	"strings"

	"github.com/daiyanuthsa/grpc-ecom-be/pb/common"
)

// BuildOrderByClause membangun klausa ORDER BY yang aman dari request pagination.
// Menerima request sort, daftar kolom yang diizinkan (whitelist), dan klausa sort default.
func BuildOrderByClause(sorts []*common.PaginationSortRequest, allowedSortFields map[string]bool, defaultSort string) (string, error) {
	if len(sorts) == 0 {
		return defaultSort, nil
	}

	var orderClauses []string
	for _, s := range sorts {
		field := strings.ToLower(s.Field)
		// Validasi field dengan whitelist
		if !allowedSortFields[field] {
			return "", fmt.Errorf("sorting by field '%s' is not allowed", s.Field)
		}

		order := strings.ToUpper(s.Order)
		// Validasi order direction
		if order != "ASC" && order != "DESC" {
			order = "ASC" // Default ke ASC jika value tidak valid
		}

		orderClauses = append(orderClauses, fmt.Sprintf(`"%s" %s`, field, order))
	}

	return "ORDER BY " + strings.Join(orderClauses, ", "), nil
}
