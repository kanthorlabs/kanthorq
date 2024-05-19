package core

import "fmt"

func QueryTruncate(table string) string {
	return fmt.Sprintf("TRUNCATE TABLE public.%s CONTINUE IDENTITY RESTRICT;", table)
}
