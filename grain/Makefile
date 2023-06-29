
BuildFlag+= -X todo.codeVersion=$(BRANCH_NAME)_$(SHORT_SHA)
export LDFLAGS=$(BuildFlag)


proto:
	protoc -I=./ --go-new_out=paths=source_relative:./ ./*.proto

.PHONY: proto