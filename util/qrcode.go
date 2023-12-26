package util

import (
	"github.com/skip2/go-qrcode"
)

type QrCodeSize int

const (
	QrCodeSmall  QrCodeSize = 256
	QrCodeMedium QrCodeSize = 512
	QrCodeBig    QrCodeSize = 1024
)

func CreateCode(content string, size QrCodeSize) []byte {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Highest, int(size))
	PanicOnError(err)
	return png
}
