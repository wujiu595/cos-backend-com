package util

import (
	"bytes"
	"image/jpeg"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func CreateQrCodeImage(content string, width int) (res []byte, err error) {
	qrcodeEnc, err := qr.Encode(content, qr.L, qr.Auto)
	if err != nil {
		return
	}

	qrcodeEnc, err = barcode.Scale(qrcodeEnc, width, width)
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	err = jpeg.Encode(buf, qrcodeEnc, &jpeg.Options{Quality: 95})
	if err != nil {
		return
	}

	res = buf.Bytes()
	return
}
