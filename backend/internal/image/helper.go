package image

import (
	"fmt"
	"math/rand/v2"
	"time"
)

func newS3Key(imageType string) string {
	t := time.Now()

	key := fmt.Sprintf("%d/%d/%d.%s", t.Year(), int(t.Month()), rand.Int(), imageType)

	return key
}
