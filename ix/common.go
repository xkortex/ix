package ix

import (
	"bufio"
	"github.com/xkortex/vprint"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

type StdInContainer struct {
	Stdin     string
	Has_stdin bool
}

// ScanLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker.
func GenScanCustomLines(seperators []string) func([]byte, bool) (int, []byte, error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		r := regexp.MustCompile("(" + strings.Join(seperators, "|") + ")")
		vprint.Println("regex: ", r)
		indexes := r.FindIndex(data)
		if len(indexes) >= 0 {
			i := indexes[0]
			return i + 1, data[0 : i+1], nil

		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	}
}

// Reads from an io.Reader and pipes it into a channel
func ScannerChannel(r io.Reader, c chan<- []byte, wg *sync.WaitGroup) {
	vprint.Println("new scanner channel")
	defer func() {
		vprint.Println("scanner channel done-ing")
		wg.Done()
	}()
	scanner := bufio.NewScanner(r)
	//splitter := GenScanCustomLines([]string{"\r\n", "\n", "\r"})
	scanner.Split(bufio.ScanLines)
	for {
		vprint.Print(".")
		if !scanner.Scan() {
			vprint.Println("breaking")
			break
		}
		out := scanner.Bytes()
		//vprint.Println(string(out))
		c <- out
	}
	vprint.Printf("Scanner finished: %v\n", r)
}

// Check if stdin is coming from a pipe
func HasStdinPipe() bool {
	info, err := os.Stdin.Stat()
	vprint.Println(info.Mode(), err)
	if err != nil {
		return false
	}
	if (info.Mode() & os.ModeCharDevice) != 0 {
		vprint.Println("Stdin isnt from pipe")
		return false
	}
	return true
}

func Get_stdin() (StdInContainer, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return StdInContainer{}, err
	}
	out_struct := StdInContainer{Has_stdin: false}
	if (info.Mode() & os.ModeCharDevice) != 0 {
		//fmt.Println("Stdin is from a terminal")
		return out_struct, nil
	}

	// data is being piped to Stdin
	//fmt.Println("data is being piped to Stdin")

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	// Deliberately block until EOF, streaming doesn't really make sense with this app
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	out_struct.Stdin = string(output)
	out_struct.Has_stdin = true
	return out_struct, nil
}
