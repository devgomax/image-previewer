package resizing

import (
	"context"
	"image"
	"net/http"

	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

// Resizer is a utility for resizing images fetched from URLs. It uses the `resize` package to perform the resizing.
type Resizer struct {
	client *http.Client
}

// NewResizer creates a new instance of Resizer with default HTTP client settings.
func NewResizer() *Resizer {
	return &Resizer{
		client: &http.Client{},
	}
}

// GetResizedImage fetches an image from URL and resizes it to the specified dimensions.
// It returns the resized image as well as the MIME type of the original image.
func (r *Resizer) GetResizedImage(ctx context.Context, url string, width, height uint, header http.Header) (image.Image, string, error) {
	imgReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", errors.Wrap(err, "[resizing::GetResizedImage]: can't create new request")
	}

	imgReq.Header = header

	resp, err := r.client.Do(imgReq)
	if err != nil {
		return nil, "", errors.Wrapf(err, "[resizing::GetResizedImage]: failed to make request to %v", url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", errors.Errorf("[resizing::GetResizedImage]: received status code %d for %s", resp.StatusCode, url)
	}

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "[resizing::GetResizedImage]: failed to decode image from response")
	}

	resized := resize.Resize(width, height, img, resize.Lanczos3)

	return resized, format, nil
}
