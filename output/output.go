package output

import (
	"strings"
	"bufio"
	"io"
)

func Write(resultsChannel<- chan string, out io.Writer) error {
	wr := bufio.NewWriter(out)
	str := &strings.Builder{}

	for result := range resultsChannel {
		str.WriteString(result)
		str.WriteRune('\n')
		_, err := wr.WriteString(str.String())
		if err != nil {
			wr.Flush()
			return err
		}
		str.Reset()
	}
	return wr.Flush()
}