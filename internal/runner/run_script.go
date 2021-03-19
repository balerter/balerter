package runner

import (
	"fmt"
	"net/http"
)

func (rnr *Runner) RunScript(name string, req *http.Request) error {
	ss, err := rnr.scriptsManager.Get()
	if err != nil {
		return err
	}

	for _, sc := range ss {
		if sc.Name == name {
			j := newJob(sc, rnr.logger)
			err = rnr.createLuaState(j, req)
			if err != nil {
				return err
			}

			rnr.jobs <- j
			return nil
		}
	}

	return fmt.Errorf("script %s not found", name)
}
