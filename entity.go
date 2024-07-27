package kanthorq

import "fmt"

func Collection(name string) string {
	return fmt.Sprintf("kanthorq_%s", name)
}
