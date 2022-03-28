package prometheus

//func TestPrometheus_getQuery_empty_query(t *testing.T) {
//	m := &Prometheus{}
//	luaState := lua.NewState()
//	_, err := m.getQuery(luaState)
//	require.Error(t, err)
//	assert.Equal(t, "query must be a string", err.Error())
//}
//
//func TestPrometheus_getQuery_query_not_a_string(t *testing.T) {
//	m := &Prometheus{}
//	luaState := lua.NewState()
//	luaState.Push(lua.LNumber(42))
//	_, err := m.getQuery(luaState)
//	require.Error(t, err)
//	assert.Equal(t, "query must be a string", err.Error())
//}
//
//func TestPrometheus_getQuery_empty_string(t *testing.T) {
//	m := &Prometheus{}
//	luaState := lua.NewState()
//	luaState.Push(lua.LString(""))
//	_, err := m.getQuery(luaState)
//	require.Error(t, err)
//	assert.Equal(t, "query must be not empty", err.Error())
//}
//
//func TestPrometheus_getQuery(t *testing.T) {
//	m := &Prometheus{}
//	luaState := lua.NewState()
//	luaState.Push(lua.LString("123"))
//	v, err := m.getQuery(luaState)
//	require.NoError(t, err)
//	assert.Equal(t, "123", v)
//}
//
//func TestPrometheus_doQuery_error_get_query(t *testing.T) {
//	m := &Prometheus{}
//	luaState := lua.NewState()
//	luaState.Push(lua.LString(""))
//	n := m.doQuery(luaState)
//	assert.Equal(t, 2, n)
//	e := luaState.Get(3)
//	assert.Equal(t, "query must be not empty", e.String())
//}
//
//func TestPrometheus_doQuery_error_parse_query_options(t *testing.T) {
//	m := &Prometheus{
//		logger: zap.NewNop(),
//	}
//
//	luaState := lua.NewState()
//	luaState.Push(lua.LString("query"))
//	tbl := &lua.LTable{}
//	tbl.RawSetString("time", lua.LNumber(42))
//	luaState.Push(tbl)
//
//	n := m.doQuery(luaState)
//	assert.Equal(t, 2, n)
//	e := luaState.Get(4)
//	assert.Equal(t, "time must be a string", e.String())
//}
//
//func TestPrometheus_doQuery_error_send(t *testing.T) {
//	hm := &httpClientMock{}
//	hm.On("Do", mock.Anything).Return(nil, fmt.Errorf("err1"))
//
//	m := &Prometheus{
//		logger: zap.NewNop(),
//		client: hm,
//		url:    &url.URL{},
//	}
//
//	luaState := lua.NewState()
//	luaState.Push(lua.LString("query"))
//
//	n := m.doQuery(luaState)
//	assert.Equal(t, 2, n)
//	e := luaState.Get(3)
//	assert.Equal(t, "error send query to prometheus: err1", e.String())
//}
//
//func TestPrometheus_doQuery_send(t *testing.T) {
//	hm := &httpClientMock{}
//	hm.On("Do", mock.Anything).Return(&http.Response{
//		Status:     "status1",
//		StatusCode: 200,
//		Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"data":{"resultType":"vector","result":[]}}`))),
//	}, nil)
//
//	m := &Prometheus{
//		logger: zap.NewNop(),
//		client: hm,
//		url:    &url.URL{},
//	}
//
//	luaState := lua.NewState()
//	luaState.Push(lua.LString("query"))
//
//	n := m.doQuery(luaState)
//	assert.Equal(t, 2, n)
//	tbl := luaState.Get(2)
//	assert.Equal(t, lua.LTTable, tbl.Type())
//}
//
//func TestPrometheus_doRange_error_get_query(t *testing.T) {
//	m := &Prometheus{}
//	luaState := lua.NewState()
//	luaState.Push(lua.LString(""))
//	n := m.doRange(luaState)
//	assert.Equal(t, 2, n)
//	e := luaState.Get(3)
//	assert.Equal(t, "query must be not empty", e.String())
//}
//
//func TestPrometheus_doRange_error_parse_query_options(t *testing.T) {
//	m := &Prometheus{
//		logger: zap.NewNop(),
//	}
//
//	luaState := lua.NewState()
//	luaState.Push(lua.LString("query"))
//	tbl := &lua.LTable{}
//	tbl.RawSetString("start", lua.LNumber(42))
//	luaState.Push(tbl)
//
//	n := m.doRange(luaState)
//	assert.Equal(t, 2, n)
//	e := luaState.Get(4)
//	assert.Equal(t, "error decode query range options, start must be a string", e.String())
//}
//
//func TestPrometheus_doRange_error_send(t *testing.T) {
//	hm := &httpClientMock{}
//	hm.On("Do", mock.Anything).Return(nil, fmt.Errorf("err1"))
//
//	m := &Prometheus{
//		logger: zap.NewNop(),
//		client: hm,
//		url:    &url.URL{},
//	}
//
//	luaState := lua.NewState()
//	luaState.Push(lua.LString("query"))
//
//	n := m.doRange(luaState)
//	assert.Equal(t, 2, n)
//	e := luaState.Get(3)
//	assert.Equal(t, "error send query to prometheus: err1", e.String())
//}
//
//func TestPrometheus_doRange_send(t *testing.T) {
//	hm := &httpClientMock{}
//	hm.On("Do", mock.Anything).Return(&http.Response{
//		Status:     "status1",
//		StatusCode: 200,
//		Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"data":{"resultType":"vector","result":[]}}`))),
//	}, nil)
//
//	m := &Prometheus{
//		logger: zap.NewNop(),
//		client: hm,
//		url:    &url.URL{},
//	}
//
//	luaState := lua.NewState()
//	luaState.Push(lua.LString("query"))
//
//	n := m.doRange(luaState)
//	assert.Equal(t, 2, n)
//	tbl := luaState.Get(2)
//	assert.Equal(t, lua.LTTable, tbl.Type())
//}
//
//func Test_processValVectorRange(t *testing.T) {
//	m := prometheus_models.Vector{}
//	m = append(m, &prometheus_models.Sample{
//		Metric:    prometheus_models.Metric{"a": "b"},
//		Value:     1,
//		Timestamp: 2,
//	})
//
//	tbl := processValVectorRange(m)
//	assert.Equal(t, lua.LTTable, tbl.Type())
//	row := tbl.RawGetInt(1)
//	require.Equal(t, lua.LTTable, row.Type())
//	metrics := row.(*lua.LTable).RawGetString("metrics")
//	require.Equal(t, lua.LTTable, metrics.Type())
//	assert.Equal(t, "b", metrics.(*lua.LTable).RawGetString("a").String())
//	assert.Equal(t, "1", row.(*lua.LTable).RawGetString("value").String())
//}
//
//func Test_processValMatrixRange(t *testing.T) {
//	m := prometheus_models.Matrix{}
//	m = append(m, &prometheus_models.SampleStream{
//		Metric: prometheus_models.Metric{"a": "b"},
//		Values: []prometheus_models.SamplePair{
//			{
//				Timestamp: 0,
//				Value:     2,
//			},
//		},
//	})
//	tbl := processValMatrixRange(m)
//	assert.Equal(t, lua.LTTable, tbl.Type())
//	row := tbl.RawGetInt(1)
//	require.Equal(t, lua.LTTable, row.Type())
//	metrics := row.(*lua.LTable).RawGetString("metrics")
//	require.Equal(t, lua.LTTable, metrics.Type())
//	assert.Equal(t, "b", metrics.(*lua.LTable).RawGetString("a").String())
//
//	vv := row.(*lua.LTable).RawGetString("values")
//	assert.Equal(t, "2", vv.(*lua.LTable).RawGetInt(1).(*lua.LTable).RawGetString("value").String())
//}
//
//func TestQueryResult_UnmarshalJSON_empty(t *testing.T) {
//	r := queryResult{}
//
//	err := r.UnmarshalJSON([]byte(``))
//	require.Error(t, err)
//	assert.Equal(t, "unexpected end of JSON input", err.Error())
//}
//
////func TestQueryResult_UnmarshalJSON_unexpected_type(t *testing.T) {
////	r := queryResult{}
////
////	err := r.UnmarshalJSON([]byte(`{"type":"foo"}`))
////	require.Error(t, err)
////	assert.Equal(t, "unexpected value type \"<ValNone>\"", err.Error())
////}
////
////func TestQueryResult_UnmarshalJSON_scalar(t *testing.T) {
////	r := queryResult{}
////
////	err := r.UnmarshalJSON([]byte(`{"resultType":"scalar","result":[1,"2"]}`))
////	require.NoError(t, err)
////	assert.Equal(t, "scalar: 2 @[1]", r.v.String())
////}
////
////func TestQueryResult_UnmarshalJSON_vector(t *testing.T) {
////	r := queryResult{}
////
////	err := r.UnmarshalJSON([]byte(`{"resultType":"vector","result":[{"metric":{},"value":[1,"2"]}]}`))
////	require.NoError(t, err)
////	assert.Equal(t, "{} => 2 @[1]", r.v.String())
////}
////
////func TestQueryResult_UnmarshalJSON_matrix(t *testing.T) {
////	r := queryResult{}
////
////	err := r.UnmarshalJSON([]byte(`{"resultType":"matrix","result":[{"metric":{},"values":[[1,"2"]]}]}`))
////	require.NoError(t, err)
////	assert.Equal(t, "{} =>\n2 @[1]", r.v.String())
////}
