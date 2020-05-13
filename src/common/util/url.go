package util

func EncodeURI(raw string) (encoded string) {
	bs := make([]byte, 0, len(raw)/2*3)
	for i := 0; i < len(raw); i++ {
		c := raw[i]
		switch c {
		// ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!#$&'()*+,-./:;=?@_~
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '!', '#', '$', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '=', '?', '@', '_', '~':
			bs = append(bs, c)
			continue
		default:
			bs = append(bs, '%')
			bs = append(bs, "0123456789ABCDEF"[c>>4])
			bs = append(bs, "0123456789ABCDEF"[c&15])
		}
	}
	encoded = string(bs)
	return
}
