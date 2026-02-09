package ncp

// RootAccount represents a root account parsed from the Excel file.
type RootAccount struct {
	AccountName string
	AccessKey   string
	SecretKey   string
	IamUsername string
	Password    string
}

// SubAccount represents a sub account returned from the NCP API.
type SubAccount struct {
	SubAccountId        string `json:"subAccountId"`
	LoginId             string `json:"loginId"`
	Name                string `json:"name"`
	Email               string `json:"email"`
	Active              bool   `json:"active"`
	Memo                string `json:"memo"`
	CanConsoleAccess    bool   `json:"canConsoleAccess"`
	CanAPIGatewayAccess bool   `json:"canApiGatewayAccess"`
	IsMfaMandatory      bool   `json:"isMfaMandatory"`
	CreateTime          string `json:"createTime"`
	UpdateTime          string `json:"updateTime"`
}

// SubAccountListResponse is the response from GET /api/v1/sub-accounts.
type SubAccountListResponse struct {
	Items      []SubAccount `json:"items"`
	TotalItems int          `json:"totalItems"`
	TotalPages int          `json:"totalPages"`
	Page       int          `json:"page"`
}

// SubAccountUpdateRequest is the request body for PUT /api/v1/sub-accounts/{id}.
type SubAccountUpdateRequest struct {
	Name               *string `json:"name,omitempty"`
	Email              *string `json:"email,omitempty"`
	Memo               *string `json:"memo,omitempty"`
	Active             *bool   `json:"active,omitempty"`
	IsMfaMandatory     *bool   `json:"isMfaMandatory,omitempty"`
	CanConsoleAccess   *bool   `json:"canConsoleAccess,omitempty"`
	CanAPIGatewayAccess *bool  `json:"canAPIGatewayAccess,omitempty"`
	UseConsolePermitIp *bool   `json:"useConsolePermitIp,omitempty"`
	UseApiAllowSource  *bool   `json:"useApiAllowSource,omitempty"`
}

// SubAccountUpdateResponse is the response from PUT /api/v1/sub-accounts/{id}.
type SubAccountUpdateResponse struct {
	Success bool `json:"success"`
}

// PasswordResetRequest is the request body for PUT /api/v1/sub-accounts/{id}/password.
type PasswordResetRequest struct {
	NeedPasswordGenerate bool    `json:"needPasswordGenerate"`
	NewPassword          *string `json:"newPassword,omitempty"`
}

// PasswordResetResponse is the response from PUT /api/v1/sub-accounts/{id}/password.
type PasswordResetResponse struct {
	Success           bool   `json:"success"`
	GeneratedPassword string `json:"generatedPassword"`
}
