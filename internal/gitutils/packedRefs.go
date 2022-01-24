package gitutils

func forEachPackedRef(packedRefsData []byte, fn func(hash string, ref string) bool) {
	offset := 0
	for offset < len(packedRefsData) {
		hash, ref, consumed := readPackedRefLine(packedRefsData, offset)
		offset += consumed

		if consumed == 0 {
			panic("Consumed 0 bytes")
		}

		if hash != "" && ref != "" {
			if !fn(hash, ref) {
				return
			}
		}
	}
}

func readPackedRefLine(packedRefsData []byte, offset int) (hash string, ref string, consumed int) {
	hash = ""
	ref = ""

	if packedRefsData[offset] == '#' {
		// This is a comment
		for i := offset; i < len(packedRefsData); i++ {
			if packedRefsData[i] == '\n' {
				consumed = i - offset + 1
				return "", "", consumed
			}
		}
		return "", "", len(packedRefsData) - offset
	}

	refStart := 0
	for i := offset; i < len(packedRefsData); i++ {
		if packedRefsData[i] == ' ' && hash == "" {
			hash = string(packedRefsData[offset:i])
			refStart = i + 1
		}
		if packedRefsData[i] == '\n' {
			ref = string(packedRefsData[refStart:i])
			consumed = i - offset + 1
			return hash, ref, consumed
		}
	}
	return

}
