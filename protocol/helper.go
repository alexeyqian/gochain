package protocol

// convert string to bytes data, leave the rest bytes empty
func CommandToBytes(command string) []byte {
	var data [commandLength]byte

	for i, c := range command {
		data[i] = byte(c)
	}

	return data[:]
}

// extract command string from bytes data
func BytesToCommand(data []byte) string {
	var command []byte

	for _, b := range data {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return string(command)
}

func ExtractCommand(request []byte) []byte {
	return request[:commandLength]
}
