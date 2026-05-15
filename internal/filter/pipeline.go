package filter

import "bufio"

// Pipeline reads lines from a Scanner, applies the Filter, and writes
// matching lines to the provided writer function.
//
// It returns the number of lines matched and any scanner error.
func Pipeline(scanner *bufio.Scanner, f *Filter, emit func(string) error) (int, error) {
	matched := 0
	for scanner.Scan() {
		line := scanner.Text()
		if f.Match(line) {
			if err := emit(line); err != nil {
				return matched, err
			}
			matched++
		}
	}
	if err := scanner.Err(); err != nil {
		return matched, err
	}
	return matched, nil
}
