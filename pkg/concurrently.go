package concurrently

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"golang.org/x/sync/errgroup"
)

type PrefixedWriter struct {
	Buf    []byte
	Mutex  *sync.Mutex
	Prefix []byte
	Writer io.Writer
}

func (w *PrefixedWriter) Flush() (int, error) {
	if len(w.Buf) == 0 {
		return 0, nil
	}

	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	_, err := w.Writer.Write(w.Prefix)
	if err != nil {
		return 0, err
	}
	count, err := w.Writer.Write(w.Buf)
	w.Buf = []byte{}

	return count, err
}

func (w *PrefixedWriter) Write(data []byte) (int, error) {
	count := 0
	for {
		index := bytes.IndexByte(data, '\n')
		if index == -1 {
			w.Buf = data
			break
		}
		w.Buf = append(w.Buf, data[:index+1]...)

		written, err := w.Flush()
		count += written
		if err != nil {
			return count, err
		}

		data = data[index+1:]
	}
	return count, nil
}

func Run(ctx context.Context, cmdSlices ...[]string) error {
	group, sctx := errgroup.WithContext(ctx)
	mutex := &sync.Mutex{}
	for index, cmdSlice := range cmdSlices {
		cmd := exec.CommandContext(sctx, cmdSlice[0], cmdSlice[1:]...)
		stdout := &PrefixedWriter{
			Buf:    []byte{},
			Mutex:  mutex,
			Prefix: fmt.Appendf([]byte{}, "[%d] ", index),
			Writer: os.Stdout,
		}
		stderr := &PrefixedWriter{
			Buf:    []byte{},
			Mutex:  mutex,
			Prefix: fmt.Appendf([]byte{}, "[%d] ", index),
			Writer: os.Stderr,
		}
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		group.Go(func() error {
			err := cmd.Run()
			stdout.Flush()
			stderr.Flush()
			return err
		})
	}

	return group.Wait()
}
