package testify

import (
	"github.com/jaswdr/faker/v2"
)

var Fake faker.Faker

func init() {
	Fake = faker.New()
}
