package model

// ModuleResponse — api_planning §3.9
type ModuleResponse struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// UpdateModuleRequest — PUT /modules/:key
type UpdateModuleRequest struct {
	Enabled *bool `json:"enabled" validate:"required"`
}
