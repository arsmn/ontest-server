package avatar

import (
	"bytes"
	"fmt"
	"image"
	"net/http"
	"strings"

	// jpeg format registration
	"image/color/palette"
	_ "image/jpeg"

	"github.com/arsmn/ontest-server/module/errors"
	"github.com/arsmn/ontest-server/module/generate"
	"github.com/disintegration/imaging"
	"github.com/issue9/identicon"
)

const AvatarSize = 290

func GenerateRandomSize(size int, data []byte) (image.Image, error) {
	randExtent := len(palette.WebSafe) - 32
	integer, err := generate.RandomInt(int64(randExtent))
	if err != nil {
		return nil, fmt.Errorf("util.RandomInt: %v", err)
	}
	colorIndex := int(integer)
	backColorIndex := colorIndex - 1
	if backColorIndex < 0 {
		backColorIndex = randExtent - 1
	}

	imgMaker, err := identicon.New(size,
		palette.WebSafe[backColorIndex], palette.WebSafe[colorIndex:colorIndex+32]...)
	if err != nil {
		return nil, fmt.Errorf("identicon.New: %v", err)
	}
	return imgMaker.Make(data), nil
}

func GenerateRandom(data []byte) (image.Image, error) {
	return GenerateRandomSize(AvatarSize, data)
}

func IsImage(data []byte) bool {
	return strings.Contains(http.DetectContentType(data), "image/")
}

func PrepareAvatar(data []byte, maxWidth, maxHeight, size int) (image.Image, error) {
	if !IsImage(data) {
		return nil, errors.ErrBadRequest.WithError("avatar is invalid")
	}

	imgCfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if imgCfg.Width > maxWidth {
		return nil, errors.ErrBadRequest.
			WithError(fmt.Sprintf("image width is too large: %d > %d", imgCfg.Width, maxWidth))
	}
	if imgCfg.Height > maxHeight {
		return nil, errors.ErrBadRequest.
			WithError(fmt.Sprintf("image height is too large: %d > %d", imgCfg.Height, maxHeight))
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	fill := imaging.Fill(img, size, size, imaging.Center, imaging.Lanczos)
	return fill, nil
}
