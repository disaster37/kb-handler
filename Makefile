.PHONY: mock-gen
mock-gen:
	go install go.uber.org/mock/mockgen@v0.3.0
	mockgen --build_flags=--mod=mod -destination=mocks/kibana_handler.go -package=mocks github.com/disaster37/kb-handler/v8 KibanaHandler