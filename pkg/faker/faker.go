package faker

import (
	lib "github.com/jaswdr/faker/v2"
)

var F lib.Faker

func init() {
	F = lib.New()
}
