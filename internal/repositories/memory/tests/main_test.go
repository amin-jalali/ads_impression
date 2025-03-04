package tests

//func TestMainSuccess(t *testing.T) {
//	os.Setenv("TEST_MODE", "true")
//	defer os.Unsetenv("TEST_MODE")
//
//	var logOutput bytes.Buffer
//	log.SetOutput(&logOutput)
//	defer log.SetOutput(os.Stderr)
//
//	oldListenAndServe := server.ListenAndServe
//	defer func() { server.ListenAndServe = oldListenAndServe }()
//
//	server.ListenAndServe = func(addr string, handler http.Handler) error {
//		return nil
//	}
//
//	var s = &http.Server{Addr: ":8080"}
//	err := s.ListenAndServe()
//	if err != nil {
//		return
//	}
//	if !bytes.Contains(logOutput.Bytes(), []byte("Server started")) {
//		t.Errorf("expected log output to contain 'Server started', but got: %s", logOutput.String())
//	}
//}
//
//func TestMainFailure(t *testing.T) {
//	os.Setenv("TEST_MODE", "true")
//	defer os.Unsetenv("TEST_MODE")
//
//	var logOutput bytes.Buffer
//	log.SetOutput(&logOutput)
//	defer log.SetOutput(os.Stderr)
//
//	oldListenAndServe := server.ListenAndServe
//	defer func() { server.ListenAndServe = oldListenAndServe }()
//
//	server.ListenAndServe = func(addr string, handler http.Handler) error {
//		return errors.New("mock ListenAndServe error")
//	}
//
//	var s = &http.Server{Addr: ":8080"}
//	err := s.ListenAndServe()
//	if err != nil {
//		return
//	}
//
//	if !bytes.Contains(logOutput.Bytes(), []byte("Server failed: mock ListenAndServe error")) {
//		t.Errorf("expected log output to contain 'Server failed', but got: %s", logOutput.String())
//	}
//}
