package gitutils

import "io"

func countLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0

	for {
		n, err := r.Read(buf)
		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				count++
			}
		}

		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
	}

	return count, nil
}
