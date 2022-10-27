package flower

import (
	"os"
	"path"
	"regexp"
)

func FindYamlFiles(dir string) ([]string, error) {
	var files []string
	dirFiles, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, f := range dirFiles {
		if f.IsDir() {
			nestedFiles, err := FindYamlFiles(path.Join(dir, f.Name()))
			if err != nil {
				return files, err
			}

			files = append(files, nestedFiles...)
			continue
		}

		match, err := regexp.MatchString("\\.ya?ml", f.Name())
		if err != nil {
			return files, err
		}
		if !match {
			continue
		}

		files = append(files, path.Join(dir, f.Name()))
	}
	return files, nil
}
