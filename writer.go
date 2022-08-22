package zlog

import (
	"bufio"
	"io"
	"runtime"
)

// SafeWriter
// @Description: 获取一个io.Writer，方便集成
// @return *io.PipeWriter
func SafeWriter() *io.PipeWriter {
	reader, writer := io.Pipe()
	go scan(reader)
	runtime.SetFinalizer(writer, writerFinalizer)
	return writer
}

func scan(reader *io.PipeReader) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanLinesOrGiveLong)
	for scanner.Scan() {
		Infof(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		Errorf("Error while reading from Writer: %s", err)
	}
	_ = reader.Close()
}

const maxTokenLength = bufio.MaxScanTokenSize / 2

func scanLinesOrGiveLong(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	if advance > 0 || token != nil || err != nil {
		return
	}
	if len(data) < maxTokenLength {
		return
	}
	return maxTokenLength, data[0:maxTokenLength], nil
}

func writerFinalizer(writer *io.PipeWriter) {
	_ = writer.Close()
}
