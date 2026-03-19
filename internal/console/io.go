package console

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

type IO interface {
	ReadLine(prompt string) (string, error)
	ReadSecret(prompt string) (string, error)
	Println(text string)
	ClearScreen()
}

type StdIO struct {
	reader    *bufio.Reader
	writer    io.Writer
	inputFile *os.File
}

func NewStdIO(reader io.Reader, writer io.Writer) *StdIO {
	inputFile, _ := reader.(*os.File)
	return &StdIO{
		reader:    bufio.NewReader(reader),
		writer:    writer,
		inputFile: inputFile,
	}
}

func (s *StdIO) ReadLine(prompt string) (string, error) {
	if _, err := fmt.Fprint(s.writer, prompt); err != nil {
		return "", err
	}
	line, err := s.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF && len(line) > 0 {
			return trimLine(line), nil
		}
		return "", err
	}
	return trimLine(line), nil
}

func (s *StdIO) ReadSecret(prompt string) (string, error) {
	if s.inputFile == nil || !term.IsTerminal(int(s.inputFile.Fd())) {
		return s.ReadLine(prompt)
	}

	if _, err := fmt.Fprint(s.writer, prompt); err != nil {
		return "", err
	}
	line, err := term.ReadPassword(int(s.inputFile.Fd()))
	if _, printErr := fmt.Fprintln(s.writer); err == nil && printErr != nil {
		return "", printErr
	}
	if err != nil {
		return "", err
	}
	return trimLine(string(line)), nil
}

func (s *StdIO) Println(text string) {
	_, _ = fmt.Fprintln(s.writer, text)
}

func (s *StdIO) ClearScreen() {
	_, _ = fmt.Fprint(s.writer, "\033[H\033[2J")
}

func trimLine(line string) string {
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1]
	}
	return line
}
