package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthZMigration struct {
	Migrations []RolePermItem `json:"authz"`
}

type SupertokensRolePermItem struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

type RolePermItem struct {
	Action string `json:"action"`
	*SupertokensRolePermItem
}

func UpdateRolePerms(jsondata string, authServerUrl string) bool {
	// TODO: use env vars
	var addUrl = fmt.Sprintf("%s/%s", authServerUrl, "add-role-perm")
	var delUrl = fmt.Sprintf("%s/%s", authServerUrl, "remove-role-perm")

	// parse
	var authz AuthZMigration
	jsonErr := json.Unmarshal([]byte(jsondata), &authz)
	if jsonErr != nil {
		fmt.Printf("UpdateRolePerms - unmarshalling: %v\n", jsonErr)
		return false
	}

	var ok bool = true
	for _, authMig := range authz.Migrations {
		endpoint := addUrl
		httpMethod := http.MethodPut

		if authMig.Action != "add" {
			// remove
			endpoint = delUrl
			httpMethod = http.MethodPost
		}

		fmt.Printf("UpdateRolePerms: %s | %s | %s\n", authMig.Action, authMig.Role, authMig.Permissions)

		// create object of type SupertokensRolePermItem
		stRoleItem := SupertokensRolePermItem{
			Role:        authMig.Role,
			Permissions: authMig.Permissions,
		}

		postBody, _ := json.Marshal(stRoleItem)

		ok = sendRequest(httpMethod, endpoint, postBody)
		if !ok {
			break
		}
	}

	return ok
}

func sendRequest(method string, endpoint string, payload []byte) bool {
	req, err := http.NewRequest(method, endpoint, bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("Http request setup error: %v\n", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Http request failed: err: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	//Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil || (resp.StatusCode != 200 && resp.StatusCode != 202) {
		bs := string(bodyBytes)
		fmt.Printf("Http request failed: status: [%d] body [%s] err: %v\n", resp.StatusCode, bs, err)
		return false
	}

	return true
}
