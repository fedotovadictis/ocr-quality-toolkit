package corpus

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"strings"
)

func makeMWSID(sourceID, imagePath string) string {
	imagePath = strings.ReplaceAll(imagePath, `\`, `/`)
	normalizedPath := path.Clean(imagePath)

	sum := sha256.Sum256([]byte(normalizedPath))
	shortHash := hex.EncodeToString(sum[:])[:12]

	return fmt.Sprintf("mws-%s-%s", sourceID, shortHash)
}
