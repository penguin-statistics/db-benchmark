package utils

import (
	"github.com/dchest/uniuri"
	"github.com/penguin-statistics/db-benchmark/models"
	"math/rand"
	"time"
)

var (
	servers = []string{"CN", "US", "JP", "KR"}
	locales = []string{"zh", "en", "ja", "ko", "zh_CN", "zh_HK", "en_US", "en_JP", "ja_JP", "ko_KR"}
	preferences = []string{"1,2", "2,3", "1,3", "1", "2", "3", "4", "1,3,4", "2,3,4", "1,2,3"}
)

func RandomFrom(c []string) string {
	l := len(c)
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(l)
	return c[i]
}

func RandomIntn(i int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(i)
}

func RandomInstallation() *models.Installation {
	return &models.Installation{
		DeviceToken:       uniuri.NewLenChars(64, []byte("0123456789abcdef")),
		Server:            RandomFrom(servers),
		Locale:            RandomFrom(locales),
		ClientPreferences: RandomFrom(preferences),
		ClientVersion:     uniuri.NewLen(16),
	}
}

// SampleHit returns true if random it with possibility. possibility should 0 < possibility < 1
func SampleHit(possibility float32) bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32() <= possibility
}
