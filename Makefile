.PHONY: mock-gen
mock-gen:
	go install github.com/golang/mock/mockgen@v1.6.0
	mockgen --build_flags=--mod=mod -destination=mocks/kibana_handler.go -package=mocks github.com/disaster37/kb-handler/v8 KibanaHandler