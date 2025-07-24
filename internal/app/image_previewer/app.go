package imagepreviewer

import (
	"github.com/devgomax/image-previewer/internal/pkg/lru"
	"github.com/devgomax/image-previewer/internal/pkg/resizing"
)

// App represents the main application logic for the image previewer. It includes caching and resizing functionalities.
type App struct {
	cache   lru.ICache
	resizer *resizing.Resizer
}

// NewApp creates a new instance of the App with specified caching and resizing configurations.
func NewApp(cache lru.ICache, resizer *resizing.Resizer) *App {
	return &App{
		cache:   cache,
		resizer: resizer,
	}
}
