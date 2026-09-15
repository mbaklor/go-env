package env

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Reads `.env` files and loads them into the program's environment variables
//
// If no filename is given, defaults to looking for `./.env` in the cwd.
// If one of the files requested doesn't exist will fail before reading any file
// to keep the env clean.
// If while parsing a file an error occurs will stop parsing that file and move
// on to the next, returning all errors recieved at the end.
func LoadFiles(files ...string) error {
	if len(files) == 0 {
		files = append(files, ".env")
	}

	// want to exit early if a file is missing so we don't get partial file reads
	for _, file := range files {
		stt, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("loading env file \"%s\": %w", file, err)
		}
		if stt.IsDir() {
			return fmt.Errorf("%s is a directory", file)
		}
	}

	errs := make([]error, 0, len(files))
	// as opposed to when checking files exist, ideally we read any file that
	// doesn't error
	for _, file := range files {
		err := parseFile(file)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse env file %s: %w", file, err))
		}

	}
	return errors.Join(errs...)
}

func parseFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	i := 0
	for s.Scan() {
		i++
		line := strings.SplitN(s.Text(), "#", 2)
		env := strings.TrimSpace(line[0])
		if env == "" {
			continue
		}
		kv := strings.SplitN(env, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("\"%s\" in line %d is not a valid env pair", s.Text(), i)
		}
		key := kv[0]
		val := kv[1]
		os.Setenv(key, val)
	}
	return s.Err()
}
