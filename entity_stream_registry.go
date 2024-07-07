package kanthorq

import "fmt"

type StreamRegistry struct {
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func StreamCollection(name string) string {
	return fmt.Sprintf("kanthorq_stream_%s", name)
}
