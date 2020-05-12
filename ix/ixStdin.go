package ix

import (
	"bufio"
	"github.com/xkortex/vprint"
	"io"
	"log"
	"os"
	"sync"
)

func RunIxStdin(slicer *MultiSlice) {
	if !HasStdinPipe() {
		log.Fatal("No stdin found")
	}
	RunIxReader(os.Stdin, slicer)
}

func RunIxReader(r io.Reader, slicer *MultiSlice) {
	var wg sync.WaitGroup
	chIn := make(chan []byte)
	defer func() {
		vprint.Println("channel in closing")
		close(chIn)
	}()
	wg.Add(1)
	go ScannerChannel(r, chIn, &wg)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		vprint.Println("closing done chan")
		close(done)
	}()

	go func() {
		buf_out := bufio.NewWriter(os.Stdout)
		var chunk_in []byte
		for {
			select {
			case <-done:
				vprint.Print("Loop completed\n")
				return
			case chunk_in = <-chIn:
				//buf_out.Write([]byte("["))
				buf_out.Write(chunk_in)
				//buf_out.Write([]byte("]"))
				buf_out.Write([]byte("\n"))
				buf_out.Flush()
			}
			chunk_in = nil
		}

	}()
	wg.Wait()
}
