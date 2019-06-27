package utils

var base36 = [36]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z',
}

// Base36Encode encode in base36
func Base36Encode(value uint64) string {
	var (
		res [16]byte
		i   int
	)

	for i = 15; value != 0; i-- {
		res[i] = base36[value%36]
		value /= 36
	}

	return string(res[i+1:])
}
