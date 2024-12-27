package apistruct

type GetConfigReq struct {
	ConfigName string `json:"configName"`
}

type GetConfigListResp struct {
	Environment string   `json:"environment"`
	Version     string   `json:"version"`
	ConfigNames []string `json:"configNames"`
}

type SetConfigReq struct {
	ConfigName string `json:"configName"`
	Data       string `json:"data"`
}
