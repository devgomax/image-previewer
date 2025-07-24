package imagepreviewer

import (
	"fmt"
	"image"

	"github.com/devgomax/image-previewer/internal/pkg/lru"
)

// getCacheKeyForImage generates a cache key for an image based on its URL, width, and height.
func getCacheKeyForImage(imageURL, width, height string) lru.Key {
	return fmt.Sprintf("%v:%v:%v", imageURL, width, height)
}

// cacheValue represents the value stored in the cache. It contains the image data and its format.
type cacheValue struct {
	img    image.Image
	format string
}
