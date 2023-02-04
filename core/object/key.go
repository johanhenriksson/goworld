package object

import "strconv"

func Key(prefix string, object T) string {
	p := len(prefix)
	buffer := make([]byte, p+1, p+9)
	copy(buffer, []byte(prefix))
	buffer[p] = '-'
	dst := strconv.AppendUint(buffer, uint64(object.ID()), 16)
	return string(dst)
}
