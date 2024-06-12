package utils

import (
	"hash/fnv"
)

func AdvisoryLockHash(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}
