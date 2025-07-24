package internalhttp

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nfnt/resize"
	"github.com/rs/zerolog/log"
)

// PreviewImage handles the preview image request.
// It takes a URL parameter for the image and two additional parameters for the width and height of the preview.
func PreviewImage(w http.ResponseWriter, r *http.Request) {
	widthParam := chi.URLParam(r, "width")
	width, err := strconv.ParseUint(widthParam, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse urlparam width")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	heightParam := chi.URLParam(r, "height")
	height, err := strconv.ParseUint(heightParam, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse urlparam height")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	imageURL, err := url.Parse(chi.URLParam(r, "*"))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse imageurl")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	imgReq, err := http.NewRequestWithContext(r.Context(), http.MethodGet, imageURL.String(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request to image server")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}

	imgReq.Header = r.Header

	client := &http.Client{}
	resp, err := client.Do(imgReq)
	if err != nil {
		log.Error().Err(err).Msg("failed to make request to image server")
		http.Error(w, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode image")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}

	resized := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)

	var buf bytes.Buffer

	switch format {
	case "jpeg", "jpg":
		if err = jpeg.Encode(&buf, resized, nil); err != nil {
			log.Error().Err(err).Msg("failed to encode jpeg")
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
	case "png":
		if err = png.Encode(&buf, resized); err != nil {
			log.Error().Err(err).Msg("failed to encode png")
			http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "image/png")
	default:
		log.Error().Str("format", format).Msg("unsupported image format")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}

	if _, err = w.Write(buf.Bytes()); err != nil {
		log.Error().Err(err).Msg("failed to write response body")
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
	}
}
