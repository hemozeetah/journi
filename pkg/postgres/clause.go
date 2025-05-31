package postgres

import (
	"fmt"
	"strings"

	"github.com/hemozeetah/journi/pkg/querybuilder"
)

var directions = map[querybuilder.Direction]string{
	querybuilder.ASC:  "ASC",
	querybuilder.DESC: "DESC",
}

var operations = map[querybuilder.Operation]string{
	querybuilder.EQ:  "=",
	querybuilder.GT:  ">",
	querybuilder.GTE: ">=",
	querybuilder.LT:  "<",
	querybuilder.LTE: "<=",
}

func WhereClause(fields map[querybuilder.Field]string, constraints []querybuilder.Constraint, data map[string]any) string {
	if len(constraints) == 0 {
		return ""
	}

	wc := make([]string, len(constraints))
	for i, v := range constraints {
		fieldRef := fmt.Sprintf("%s%d", fields[v.Field], i)
		data[fieldRef] = v.Value

		wc[i] = fmt.Sprintf("%s %s :%s", fields[v.Field], operations[v.Operation], fieldRef)
	}

	return fmt.Sprintf("WHERE %s", strings.Join(wc, " AND "))
}

func OrderByClause(fields map[querybuilder.Field]string, orderBy querybuilder.OrderBy) string {
	return fmt.Sprintf("ORDER BY %s %s", fields[orderBy.Field], directions[orderBy.Direction])
}

func OffsetClause(fields map[querybuilder.Field]string, page querybuilder.Page) string {
	return fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", (page.Number-1)*page.Rows, page.Rows)
}
