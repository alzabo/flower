package flower

import (
	"os"
	"path"
	"regexp"
)

func FindYamlFiles(dirs []string) ([]string, error) {
	var files []string
	var err error
	for _, dir := range dirs {
		f, err := walkForYaml(dir)
		if err != nil {
			return files, err
		}
		files = append(files, f...)
	}
	return files, err
}

func walkForYaml(dir string) ([]string, error) {
	var files []string
	dirFiles, err := os.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, f := range dirFiles {
		if f.IsDir() {
			nestedFiles, err := walkForYaml(path.Join(dir, f.Name()))
			if err != nil {
				return files, err
			}

			files = append(files, nestedFiles...)
			continue
		}

		match, err := regexp.MatchString(`\.ya?ml$`, f.Name())
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
