package ix

import (
	"bufio"
	"bytes"
	"github.com/xkortex/vprint"
	"io"
	"log"
	"os"
	"sync"
)

func b2s(ary [][]byte, sep string ) string {
	s := ""
	for _, b := range ary {
		s += string(b) + sep
	}
	return s
}

func RunIxStdin(slicer *MultiSlice) {
	if !HasStdinPipe() {
		log.Fatal("No stdin found")
	}
	RunIxReader(os.Stdin, slicer)
}

func IxRecordSlicer(slicer *MultiSlice, chIn <-chan []byte, chOut chan<- []byte, done chan struct{}) {
	var chunk_in []byte
	fieldSep := []byte(slicer.Sep)
	vprint.Printf("fieldSep: `%v' %v\n", slicer.Sep, fieldSep)
	for {
		select {
		case <-done:
			return
		case chunk_in = <-chIn:
			vprint.Printf("chIn: `%v'\n", string(chunk_in))
			parts := bytes.Split(chunk_in, fieldSep)
			thisSlicer := slicer.FieldSlicer.Copy()
			thisSlicer.Normalize(len(parts))

			vprint.Printf("parts: %d [%v:%v]: '%v'", len(parts), thisSlicer.Start, thisSlicer.Stop, b2s(parts, "|"))
			chOut <- bytes.Join(parts[thisSlicer.Start:thisSlicer.Stop], fieldSep)
		}
	}
}

func RunIxReader(r io.Reader, slicer *MultiSlice) {
	var wg sync.WaitGroup
	chIn := make(chan []byte)
	chOut := make(chan []byte)
	defer func() {
		vprint.Println("channel in closing")
		close(chIn)
		close(chOut)
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
		go IxRecordSlicer(slicer, chIn, chOut, done)
		for {
			select {
			case <-done:
				vprint.Print("Loop completed\n")
				return
			case chunk_in = <-chOut:
				buf_out.Write([]byte("["))
				buf_out.Write(chunk_in)
				buf_out.Write([]byte("]"))
				buf_out.Write([]byte("\n"))
				buf_out.Flush()
			}
			chunk_in = nil
		}

	}()
	wg.Wait()
}
