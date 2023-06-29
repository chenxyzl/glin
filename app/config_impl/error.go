package config_impl

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"laiya/proto/code"
)

type ErrorConfig struct {
	Error map[string]string
}

func (conf *ErrorConfig) Parse(data []byte) error {
	//parse
	_, err := toml.Decode(string(data), conf)
	if err != nil {
		return err
	}
	//check
	var codeMapping = map[code.Code]string{}
	for key, val := range conf.Error {
		cod, ok := code.Code_value[key]
		if !ok {
			return fmt.Errorf("error key not in code, key:%v|val:%v", key, val)
		}
		codeMapping[code.Code(cod)] = val
	}
	//set to code
	code.SetCodeMappingMsg(codeMapping)
	return nil
}
