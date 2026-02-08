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
	SubAccountId       int64  `json:"subAccountId"`
	LoginId            string `json:"loginId"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Active             bool   `json:"active"`
	Memo               string `json:"memo"`
	CanConsoleAccess   bool   `json:"canConsoleAccess"`
	CanAPIGatewayAccess bool  `json:"canApiGatewayAccess"`
	IsMfaMandatory     bool   `json:"isMfaMandatory"`
	CreateTime         string `json:"createTime"`
	UpdateTime         string `json:"updateTime"`
}

// SubAccountListResponse is the response from GET /api/v1/sub-accounts.
type SubAccountListResponse struct {
	Content    []SubAccount `json:"content"`
	TotalRows  int          `json:"totalRows"`
	TotalPages int          `json:"totalPages"`
	Page       int          `json:"page"`
	PageSize   int          `json:"pageSize"`
}

// SubAccountUpdateRequest is the request body for PUT /api/v1/sub-accounts/{id}.
type SubAccountUpdateRequest struct {
	Email              *string `json:"email,omitempty"`
	Name               *string `json:"name,omitempty"`
	Memo               *string `json:"memo,omitempty"`
	IsMfaMandatory     *bool   `json:"isMfaMandatory,omitempty"`
	UseConsolePermitIp *bool   `json:"useConsolePermitIp,omitempty"`
	UseApiAllowSource  *bool   `json:"useApiAllowSource,omitempty"`
	Active             *bool   `json:"active,omitempty"`
	Password           *string `json:"password,omitempty"`
	NeedPasswordReset  *bool   `json:"needPasswordReset,omitempty"`
}
