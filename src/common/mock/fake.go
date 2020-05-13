package mock

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"cos-backend-com/src/common/flake"
)

func FakeMacaddr(name string) string {
	h := md5.New()
	h.Write([]byte(name))
	encoded := hex.EncodeToString(h.Sum(nil))[:12]
	return encoded[0:2] + ":" + encoded[2:4] + ":" + encoded[4:6] +
		":" + encoded[6:8] + ":" + encoded[8:10] + ":" + encoded[10:12]
}

func FakeId(name string) flake.ID {
	h := md5.New()

	var hasSeq *string
	n := strings.LastIndex(name, "-")
	if n != -1 {
		v, err := strconv.ParseInt(name[n+1:], 10, 64)
		if v > 0xffff {
			panic("seqid overflow 0xffff")
		}
		if err == nil {
			seq := strconv.FormatInt(v, 16)
			hasSeq = &seq
			name = name[:n]
		}
	}

	h.Write([]byte(name))
	encoded := hex.EncodeToString(h.Sum(nil))[:10]

	if hasSeq != nil {
		encoded += fmt.Sprintf("%04s", *hasSeq)
	}

	v, _ := strconv.ParseInt(encoded, 16, 64)
	return flake.ID(v)
}
