package imagepreviewer

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// PreviewImage handles the preview image request.
// It takes a URL parameter for the image and two additional parameters for the width and height of the preview.
func (a *App) PreviewImage(w http.ResponseWriter, r *http.Request) {
	var (
		format  string
		resized image.Image
		err     error
	)

	widthParam := chi.URLParam(r, "width")
	width, err := strconv.ParseUint(widthParam, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("[image_previewer::PreviewImage]: failed to parse urlparam width")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	heightParam := chi.URLParam(r, "height")
	height, err := strconv.ParseUint(heightParam, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("[image_previewer::PreviewImage]: failed to parse urlparam height")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	imageURL, err := url.Parse(chi.URLParam(r, "*"))
	if err != nil {
		log.Error().Err(err).Msg("[image_previewer::PreviewImage]: failed to parse imageurl")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	val, ok := a.cache.Get(getCacheKeyForImage(imageURL.String(), widthParam, heightParam))
	if ok {
		cacheVal := val.(cacheValue)
		format, resized = cacheVal.format, cacheVal.img
	} else {
		resized, format, err = a.resizer.GetResizedImage(r.Context(), imageURL.String(), uint(width), uint(height), r.Header)
		if err != nil {
			log.Error().Err(err).Msg("[image_previewer::PreviewImage]: failed to get resized image")
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}

		a.cache.Set(getCacheKeyForImage(imageURL.String(), widthParam, heightParam), cacheValue{
			img:    resized,
			format: format,
		})
	}

	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		if err = jpeg.Encode(&buf, resized, nil); err != nil {
			log.Error().Err(err).Msg("[image_previewer::PreviewImage]: failed to encode jpeg")
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
	case "png":
		if err = png.Encode(&buf, resized); err != nil {
			log.Error().Err(err).Msg("[image_previewer::PreviewImage]: failed to encode png")
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "image/png")
	default:
		log.Error().Str("format", format).Msg("[image_previewer::PreviewImage]: unsupported image format")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}

	if _, err = w.Write(buf.Bytes()); err != nil {
		log.Error().Err(err).Msg("[image_previewer::PreviewImage]: failed to write response body")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
}
