package alerts

import (
	"fmt"
	alert2 "github.com/balerter/balerter/internal/alert"
	"net/http"
	"strings"
)

func parseNames(argValue string) map[string]struct{} {
	result := map[string]struct{}{}
	if argValue == "" {
		return result
	}

	names := strings.Split(argValue, ",")
	for _, n := range names {
		result[n] = struct{}{}
	}

	return result
}

func parseLevels(argValue string) (map[alert2.Level]struct{}, error) {
	result := map[alert2.Level]struct{}{}
	if argValue == "" {
		return result, nil
	}

	levels := strings.Split(argValue, ",")
	for _, l := range levels {
		ll, err := alert2.LevelFromString(l)
		if err != nil {
			return nil, fmt.Errorf("bad level value")
		}

		result[ll] = struct{}{}
	}

	return result, nil
}

func filter(req *http.Request, data []*alert2.Alert) ([]*alert2.Alert, error) {
	levelsMap, err := parseLevels(req.URL.Query().Get("level"))
	if err != nil {
		return nil, err
	}

	namesMap := parseNames(req.URL.Query().Get("name"))

	var result []*alert2.Alert

	for _, item := range data {
		if _, ok := levelsMap[item.Level]; len(levelsMap) > 0 && !ok {
			continue
		}
		if _, ok := namesMap[item.Name]; len(namesMap) > 0 && !ok {
			continue
		}

		result = append(result, item)
	}

	return result, nil
}
