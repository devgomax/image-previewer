package e2e

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	_ "image/jpeg"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	client      = http.DefaultClient
	urlTemplate = "http://localhost:8081/fill/%v/%v/%v"
	imgTemplate = "http://nginx:80/%v"
)

//go:embed data/gopher_2000x1000.jpg
var gopher2000x1000 []byte

//go:embed data/gopher_3000x1500.jpg
var gopher3000x1500 []byte

//go:embed data/gopher_1000x500.jpg
var gopher1000x500 []byte

//go:embed data/gopher_200x200.jpg
var gopher200x200 []byte

func TestGetExistingImage(t *testing.T) {
	tests := []struct {
		imageName string
		file      []byte
		width     int
		height    int
	}{
		{
			imageName: "gopher_200x200.jpg",
			file:      gopher200x200,
			width:     200,
			height:    200,
		},
		{
			imageName: "gopher_1000x500.jpg",
			file:      gopher1000x500,
			width:     1000,
			height:    500,
		},
		{
			imageName: "gopher_2000x1000.jpg",
			file:      gopher2000x1000,
			width:     2000,
			height:    1000,
		},
		{
			imageName: "gopher_3000x1500.jpg",
			file:      gopher3000x1500,
			width:     3000,
			height:    1500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.imageName, func(t *testing.T) {
			imgURL := fmt.Sprintf(imgTemplate, "gopher_2000x1000.jpg")
			reqURL := fmt.Sprintf(urlTemplate, tt.width, tt.height, imgURL)

			req, err := http.NewRequest(http.MethodGet, reqURL, nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			img, format, err := image.Decode(resp.Body)
			require.NoError(t, err)

			expectedImg, expectedFormat, err := image.Decode(bytes.NewBuffer(tt.file))
			require.NoError(t, err)
			require.Equal(t, expectedFormat, format)
			require.Equal(t, expectedImg.Bounds(), img.Bounds())
		})
	}
}

func TestNegative(t *testing.T) {
	tests := []struct {
		name   string
		width  any
		height any
		imgURL string
		status int
	}{
		{
			name:   "width is not a number",
			width:  "hundred",
			height: 100,
			imgURL: fmt.Sprintf(imgTemplate, "gopher_2000x1000.jpg"),
			status: http.StatusBadRequest,
		},
		{
			name:   "height is not a number",
			width:  100,
			height: "hundred",
			imgURL: fmt.Sprintf(imgTemplate, "gopher_2000x1000.jpg"),
			status: http.StatusBadRequest,
		},
		{
			name:   "width is negative number",
			width:  -100,
			height: 100,
			imgURL: fmt.Sprintf(imgTemplate, "gopher_2000x1000.jpg"),
			status: http.StatusBadRequest,
		},
		{
			name:   "height is negative number",
			width:  100,
			height: -100,
			imgURL: fmt.Sprintf(imgTemplate, "gopher_2000x1000.jpg"),
			status: http.StatusBadRequest,
		},
		{
			name:   "link is not a valid url",
			width:  100,
			height: -100,
			imgURL: "gopher_2000x1000.jpg",
			status: http.StatusBadRequest,
		},
		{
			name:   "image not found",
			width:  100,
			height: 100,
			imgURL: fmt.Sprintf(imgTemplate, "gopher_1000x1000.jpg"),
			status: http.StatusBadGateway,
		},
		{
			name:   "server not found",
			width:  100,
			height: 100,
			imgURL: "http://localhost/gopher_1000x1000.jpg",
			status: http.StatusBadGateway,
		},
		{
			name:   "not an image file",
			width:  100,
			height: 100,
			imgURL: fmt.Sprintf(imgTemplate, "gopher.txt"),
			status: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqURL := fmt.Sprintf(urlTemplate, tt.width, tt.height, tt.imgURL)

			req, err := http.NewRequest(http.MethodGet, reqURL, nil)
			require.NoError(t, err)

			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tt.status, resp.StatusCode)
		})
	}
}
