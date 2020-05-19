package provider

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/balerter/balerter/internal/script/script"
)

func ReadScript(fs Fs, scriptPath string) (*script.Script, error) {
	reader, err := fs.Open(scriptPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	_, fn := path.Split(scriptPath)
	s := script.New()
	s.Name = strings.TrimSuffix(fn, ".lua")
	s.Body = body

	if err := s.ParseMeta(); err != nil {
		return nil, err
	}

	return s, nil
}
