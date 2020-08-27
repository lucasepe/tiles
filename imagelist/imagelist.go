package imagelist

import (
	"bufio"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Load grabs the list off all the images
// in the folder (id uri is a folder) or in the
// text file (if uri starts with the '@' character).
func Load(uri string) ([]string, error) {
	if strings.HasPrefix(uri, "@") {
		filename, err := resolveTildeEventually(uri[1:])
		if err != nil {
			return nil, err
		}
		return FromFile(filename)
	}

	dirname, err := resolveTildeEventually(uri)
	if err != nil {
		return nil, err
	}

	//fmt.Fprintf(os.Stderr, "loading image list from folder <%s>\n", dirname)
	return FromFolder(dirname)
}

// FromFile returns a slice with all PNG
// images path listed in the specified text file.
func FromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	res := []string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); len(line) > 0 {
			res = append(res, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// FromFolder returns a slice with all PNG images
// located in 'dirname'.
func FromFolder(dirname string) ([]string, error) {
	fp, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	files, err := fp.Readdir(-1)
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, el := range files {
		if el.IsDir() {
			continue
		}

		if strings.EqualFold(filepath.Ext(el.Name()), ".png") {
			res = append(res, filepath.Join(dirname, el.Name()))
		}
	}

	return res, nil
}

// resolveTildeEventually expand the `~` character
// as the user home directory.
func resolveTildeEventually(uri string) (string, error) {
	if strings.HasPrefix(uri, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}

		return filepath.Join(usr.HomeDir, uri[1:]), nil
	}

	return uri, nil
}
