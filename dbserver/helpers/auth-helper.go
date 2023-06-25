package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AuthZMigration struct {
	Migrations []RolePermItem `json:"authz"`
}

type RolePermItem struct {
	Action      string   `json:"action"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

func AddRolePerms(jsondata string) error {
	// parse
	var authz AuthZMigration
	jsonErr := json.Unmarshal([]byte(jsondata), &authz)
	if jsonErr != nil {
		return jsonErr
	}

	for _, authMig := range authz.Migrations {
		if authMig.Action == "add" {
			fmt.Printf("Processing: %s, %s, %s\n", authMig.Action, authMig.Role, authMig.Permissions)
			postBody, _ := json.Marshal(jsondata)
			post(postBody)
		}
	}

	return nil
}

func post(postBody []byte) {
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post("https://localhost:7567/add-role-perm", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Println(sb)
}
