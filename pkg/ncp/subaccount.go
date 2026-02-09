package ncp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ListSubAccounts retrieves all sub accounts for the authenticated root account.
func (c *Client) ListSubAccounts() ([]SubAccount, error) {
	var all []SubAccount
	page := 0
	pageSize := 100

	for {
		path := fmt.Sprintf("/api/v1/sub-accounts?pageSize=%d&page=%d", pageSize, page)
		body, statusCode, err := c.doRequest("GET", path, nil)
		if err != nil {
			return nil, fmt.Errorf("listing sub accounts: %w", err)
		}
		if statusCode != 200 {
			return nil, fmt.Errorf("listing sub accounts: HTTP %d - %s", statusCode, string(body))
		}

		var resp SubAccountListResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, fmt.Errorf("parsing response: %w", err)
		}

		all = append(all, resp.Items...)

		if len(all) >= resp.TotalItems || len(resp.Items) == 0 {
			break
		}
		page++
	}

	return all, nil
}

// UpdateSubAccount updates a sub account by ID.
func (c *Client) UpdateSubAccount(subAccountId string, req *SubAccountUpdateRequest) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/sub-accounts/%s", subAccountId)
	body, statusCode, err := c.doRequest("PUT", path, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("updating sub account: %w", err)
	}
	if statusCode != 200 {
		return fmt.Errorf("updating sub account: HTTP %d - %s", statusCode, string(body))
	}

	var resp SubAccountUpdateResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("parsing update response: %w", err)
	}
	if !resp.Success {
		return fmt.Errorf("API returned success=false: %s", string(body))
	}

	return nil
}

// ResetPassword resets a sub account's login password.
// If password is empty, it auto-generates a password and returns it.
// If password is provided, it sets the password to the given value.
func (c *Client) ResetPassword(subAccountId string, password string) (string, error) {
	var req PasswordResetRequest
	if password == "" {
		req.NeedPasswordGenerate = true
	} else {
		req.NeedPasswordGenerate = false
		req.NewPassword = &password
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshaling password request: %w", err)
	}

	path := fmt.Sprintf("/api/v1/sub-accounts/%s/password", subAccountId)
	body, statusCode, err := c.doRequest("PUT", path, bytes.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("resetting password: %w", err)
	}
	if statusCode != 200 {
		return "", fmt.Errorf("resetting password: HTTP %d - %s", statusCode, string(body))
	}

	var resp PasswordResetResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("parsing password response: %w", err)
	}
	if !resp.Success {
		return "", fmt.Errorf("password reset failed: %s", string(body))
	}

	return resp.GeneratedPassword, nil
}

// ActivateSubAccount activates a sub account and resets its password.
// Returns the generated password (if auto-generated) and any error.
func (c *Client) ActivateSubAccount(sa SubAccount, password string) (string, error) {
	// 1. 활성화
	active := true
	req := &SubAccountUpdateRequest{
		Name:   &sa.Name,
		Active: &active,
	}
	if err := c.UpdateSubAccount(sa.SubAccountId, req); err != nil {
		return "", err
	}

	// 2. 비밀번호 재설정 (별도 API)
	generatedPw, err := c.ResetPassword(sa.SubAccountId, password)
	return generatedPw, err
}

// DeactivateSubAccount deactivates (suspends) a sub account.
func (c *Client) DeactivateSubAccount(sa SubAccount) error {
	active := false
	req := &SubAccountUpdateRequest{
		Name:   &sa.Name,
		Active: &active,
	}
	return c.UpdateSubAccount(sa.SubAccountId, req)
}
