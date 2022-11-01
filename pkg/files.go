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
	yamlFileExpr := regexp.MustCompile(`\.ya?ml$`)
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

		match := yamlFileExpr.MatchString(f.Name())
		if !match {
			continue
		}

		files = append(files, path.Join(dir, f.Name()))
	}
	return files, nil
}
