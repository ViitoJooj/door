package dtos

type Env struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UpgradeEnvInput struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EnvOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Env    `json:"data"`
}
