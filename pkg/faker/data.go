package faker

import (
	"encoding/json"
	"log"
)

func DataOfSize(size int) []byte {
	if size < 0 {
		log.Panicf("invalid size: %d", size)
	}

	// need reserve 1kb for metadata
	data := F.RandomStringWithLength(size - 1024)
	bytes, err := json.Marshal(map[string]string{"data": data})
	if err != nil {
		panic(err)
	}

	return bytes
}

func DataOf16Kb() []byte {
	return DataOfSize(16 * 1024)
}
