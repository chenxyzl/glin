package config_impl

import "encoding/json"

var inviteBroadcastGameTypesConf []string

func (conf *ClientDynamic) GetInviteBroadcastGameTypes() []string {
	return inviteBroadcastGameTypesConf
}

func (conf *ClientDynamic) parseInviteBroadcastGameTypes() error {
	var config []string
	err := json.Unmarshal([]byte((*conf)["InviteBroadcastGameTypes"]), &config)
	if err != nil {
		return err
	}
	inviteBroadcastGameTypesConf = config
	return nil
}

func (conf *ClientDynamic) GetConfig(key string) string {
	return (*conf)[key]
}
