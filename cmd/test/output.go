package main

import (
	"encoding/json"
	"fmt"
	"github.com/balerter/balerter/internal/modules"
	"github.com/fatih/color"
	"io"
)

func output(results []modules.TestResult, w io.Writer, asJson bool) error {
	if asJson {
		return outputJson(results, w)
	}

	return outputPlainColored(results, w)
}

func outputPlainColored(results []modules.TestResult, w io.Writer) error {

	colorOk := color.New(color.FgGreen)
	colorFail := color.New(color.FgRed)

	for _, r := range results {
		line := ""
		if r.Ok {
			line += colorOk.Sprint("[PASS]")
		} else {
			line += colorFail.Sprint("[FAIL]")
		}

		line += "\t[" + r.ScriptName + "]"
		line += "\t[" + r.ModuleName + "]"
		line += "\t" + r.Message

		_, err := fmt.Fprintf(w, "%s\n", line)
		if err != nil {
			return err
		}
	}

	return nil
}

func outputJson(results []modules.TestResult, w io.Writer) error {
	e, err := json.Marshal(results)
	if err != nil {
		return nil
	}

	_, err = fmt.Fprintf(w, string(e))
	return err
}
