package idx

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

// @TODO: consider usinghttps://github.com/jetify-com/typeid
func New(ns string) string {
	return fmt.Sprintf("%s_%s", ns, ksuid.New().String())
}

func Build(ns, id string) string {
	return fmt.Sprintf("%s_%s", ns, id)
}

func FromTime(ns string, t time.Time) string {
	// error could not be happen because we provide a valid payload
	id, _ := ksuid.FromParts(t, DefaultPayload)
	return fmt.Sprintf("%s_%s", ns, id.String())
}

func ToTime(nsId string) (time.Time, error) {
	segments := strings.Split(nsId, "_")
	if len(segments) != 2 {
		return time.Time{}, errors.New("IDX.MALFORMED_FORMAT.ERROR")
	}

	id, err := ksuid.Parse(segments[1])
	if err != nil {
		return time.Time{}, errors.New("IDX.PARSE.ERROR")
	}

	return id.Time(), nil
}
