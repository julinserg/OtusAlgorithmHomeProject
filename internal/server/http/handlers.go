package internalhttp

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/julinserg/OtusAlgorithmHomeProject/internal/app"
)

var (
	ErrNotSetParams          = errors.New("not set width, height or image path in URL")
	ErrWidthNotInt           = errors.New("width in URL not integer")
	ErrHeightNotInt          = errors.New("height in URL not integer")
	ErrDimmensionIsVeryLarge = errors.New("width or height is very large")
)

type minisearchHandler struct {
	logger Logger
	app    Application
}

func (ph *minisearchHandler) hellowHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("This is my minisearch!"))
}

func (ph *minisearchHandler) sendError(err error, code int, w http.ResponseWriter) {
	ph.logger.Error(err.Error())
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}

func (ph *minisearchHandler) proxyError(body []byte, code int, w http.ResponseWriter) {
	ph.logger.Error("Image Server Error Code: " + strconv.Itoa(code))
	w.WriteHeader(code)
	w.Write(body)
}

func (ph *minisearchHandler) validateInputParameter(url string, w http.ResponseWriter) (bool, app.InputParams) {
	splitURL := strings.Split(url, "/")
	if len(splitURL) < 3 {
		ph.sendError(ErrNotSetParams, http.StatusBadRequest, w)
		return false, app.InputParams{}
	}
	splitURLParam := splitURL[2:]

	if len(splitURLParam) < 3 {
		ph.sendError(ErrNotSetParams, http.StatusBadRequest, w)
		return false, app.InputParams{}
	}
	width, err := strconv.Atoi(splitURLParam[0])
	if err != nil {
		ph.sendError(ErrWidthNotInt, http.StatusBadRequest, w)
		return false, app.InputParams{}
	}

	height, err := strconv.Atoi(splitURLParam[1])
	if err != nil {
		ph.sendError(ErrHeightNotInt, http.StatusBadRequest, w)
		return false, app.InputParams{}
	}

	if width > 3840 || height > 2160 {
		ph.sendError(ErrDimmensionIsVeryLarge, http.StatusBadRequest, w)
		return false, app.InputParams{}
	}

	imageURL := strings.Join(splitURLParam[2:], "/")

	return true, app.InputParams{Width: width, Height: height, ImageURL: imageURL}
}

func (ph *minisearchHandler) mainHandler(w http.ResponseWriter, r *http.Request) {
	isValid, params := ph.validateInputParameter(r.URL.Path, w)
	if !isValid {
		return
	}
	image, code, isFromCache, err := ph.app.GetImagePreview(params, r.Header)
	if err != nil {
		if errors.Is(err, app.ErrFromRemoteServer) {
			ph.proxyError(image, code, w)
		} else {
			ph.sendError(err, code, w)
			return
		}
	}
	w.Header().Add("is-image-from-cache", strconv.FormatBool(isFromCache))
	w.Write(image)
	w.WriteHeader(http.StatusOK)
}

func (ph *minisearchHandler) clearCacheHandler(w http.ResponseWriter, r *http.Request) {
	ph.app.ClearCache()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Clear cache is done!"))
}
