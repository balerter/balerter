package alerts

import (
	"fmt"
	"github.com/balerter/balerter/internal/alert/alert"
	alertManager "github.com/balerter/balerter/internal/alert/manager"
	"net/http"
	"strings"
)

func filter(req *http.Request, data []*alertManager.AlertInfo) ([]*alertManager.AlertInfo, error) {
	levelsMap := map[alert.Level]struct{}{}
	namesMap := map[string]struct{}{}

	levels := strings.Split(req.URL.Query().Get("level"), ",")
	for _, l := range levels {
		if l == "" {
			continue
		}
		ll, err := alert.LevelFromString(l)
		if err != nil {
			return nil, fmt.Errorf("bad level value")
		}

		levelsMap[ll] = struct{}{}
	}

	names := strings.Split(req.URL.Query().Get("name"), ",")
	for _, n := range names {
		if n == "" {
			continue
		}
		namesMap[n] = struct{}{}
	}

	var result []*alertManager.AlertInfo

	for _, item := range data {
		if len(levelsMap) > 0 {
			if _, ok := levelsMap[item.Level]; !ok {
				continue
			}
		}
		if len(namesMap) > 0 {
			if _, ok := namesMap[item.Name]; !ok {
				continue
			}
		}

		result = append(result, item)
	}

	return result, nil
}
