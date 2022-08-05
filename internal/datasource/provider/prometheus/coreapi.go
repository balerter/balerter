package prometheus

import (
	"fmt"
	"net/http"
)

type runapiRequest struct {
	Query   string `json:"query"`
	Range   string `json:"range"`
	Options struct {
		Time  int `json:"time"`
		Start int `json:"start"`
		End   int `json:"end"`
		Step  int `json:"step"`
	} `json:"options"`
}

func (r runapiRequest) validate() error {
	if r.Query == "" && r.Range == "" {
		return fmt.Errorf("query or range must be set")
	}

	if r.Query != "" && r.Range != "" {
		return fmt.Errorf("query and range cannot be set at the same time")
	}

	return nil
}

func (m *Prometheus) CoreApiHandler(req []string, body []byte) (any, int, error) {
	return nil, http.StatusNotImplemented, fmt.Errorf("not implemented")
}

//r := runapiRequest{}
//
//err := json.NewDecoder(req.Body).Decode(&r)
//if err != nil {
//	http.Error(rw, "error decode request: "+err.Error(), http.StatusBadRequest)
//	return
//}
//
//if errValidate := r.validate(); errValidate != nil {
//	http.Error(rw, "error validate request: "+errValidate.Error(), http.StatusBadRequest)
//	return
//}
//
//var u string
//if r.Query != "" {
//	opts := &queryQueryOptions{}
//	if r.Options.Time != 0 {
//		opts.Time = strconv.Itoa(r.Options.Time)
//	}
//	u = m.sendQuery(r.Query, opts)
//} else {
//	opts := &queryRangeOptions{}
//	if r.Options.Start != 0 {
//		opts.Start = strconv.Itoa(r.Options.Start)
//	}
//	if r.Options.End != 0 {
//		opts.End = strconv.Itoa(r.Options.End)
//	}
//	if r.Options.Step != 0 {
//		opts.Step = strconv.Itoa(r.Options.Step)
//	}
//	u = m.sendRange(r.Range, opts)
//}
//
//res, errDo := m.send(u)
//if errDo != nil {
//	http.Error(rw, "error do query: "+errDo.Error(), http.StatusInternalServerError)
//	return
//}
//
//resp, errMarshal := json.Marshal(res)
//if errMarshal != nil {
//	http.Error(rw, "error marshal response: "+errMarshal.Error(), http.StatusInternalServerError)
//	return
//}
//
//rw.Write(resp)
