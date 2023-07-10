package pkg

import (
	"log"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/userroles"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type RolePerm struct {
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

func InitSupertokensAuth(supertokensServerUrl string) error {
	log.Println("InitSupertokensAuth : >>> starting...")
	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// https://try.supertokens.com is for demo purposes. Replace this with the address of your core instance (sign up on supertokens.com), or self host a core.
			ConnectionURI: supertokensServerUrl,
			// APIKey: <API_KEY(if configured)>,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "Unbuilt",
			APIDomain:       "http://localhost:5000",
			WebsiteDomain:   "http://localhost:3000",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			thirdparty.Init(&tpmodels.TypeInput{
				SignInAndUpFeature: tpmodels.TypeInputSignInAndUp{
					Providers: []tpmodels.TypeProvider{
						// We have provided you with development keys which you can use for testing.
						// IMPORTANT: Please replace them with your own OAuth keys for production use.
						thirdparty.Google(tpmodels.GoogleConfig{
							ClientID:     "1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com",
							ClientSecret: "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",
						}),
						thirdparty.Github(tpmodels.GithubConfig{
							ClientID:     "467101b197249757c71f",
							ClientSecret: "e97051221f4b6426e8fe8d51486396703012f5bd",
						}),
						// thirdparty.Apple(tpmodels.AppleConfig{
						// 	ClientID: "4398792-io.supertokens.example.service",
						// 	ClientSecret: tpmodels.AppleClientSecret{
						// 		KeyId:      "7M48Y4RYDL",
						// 		PrivateKey: "-----BEGIN PRIVATE KEY-----\nMIGTAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBHkwdwIBAQQgu8gXs+XYkqXD6Ala9Sf/iJXzhbwcoG5dMh1OonpdJUmgCgYIKoZIzj0DAQehRANCAASfrvlFbFCYqn3I2zeknYXLwtH30JuOKestDbSfZYxZNMqhF/OzdZFTV0zc5u5s3eN+oCWbnvl0hM+9IW0UlkdA\n-----END PRIVATE KEY-----",
						// 		TeamId:     "YWQCXGJRJL",
						// 	},
						// }),
						// thirdparty.Facebook(tpmodels.FacebookConfig{
						//    ClientID:     "FACEBOOK_CLIENT_ID",
						//    ClientSecret: "FACEBOOK_CLIENT_SECRET",
						// }),
					}}}),
			session.Init(nil),   // initializes session features
			userroles.Init(nil), // initialise userroles features
		},
	})

	if err != nil {
		log.Printf("InitSupertokensAuth : could not initialise! Err: %s", err.Error())
		return err
	}

	log.Println("InitSupertokensAuth : <<< Done.")
	return nil
}

// ---------------------------------------------------------------------
//
//	All supertokens-facing methods begin with ST*
//
// --------------------------------------------------------------------
func STAddRolePerm(role string, perms []string) error {
	/**
	* You can choose to give multiple or no permissions when creating a role
	* createNewRoleOrAddPermissions("user", []string{}) - No permissions
	* createNewRoleOrAddPermissions("user", []string{"read", "write"}) - Multiple permissions
	 */
	resp, err := userroles.CreateNewRoleOrAddPermissions(role, perms, nil)

	if err != nil {
		log.Printf("CreateNewRoleOrAddPermissions failed! %v", err)
		return err
	}

	if !resp.OK.CreatedNewRole {
		// The role already exists
		log.Printf("⚠️ [%s|%v] already exists.", role, perms)
	}
	return nil
}

func STDelRolePerm(role string, perms []string) error {
	resp, err := userroles.RemovePermissionsFromRole(role, perms, nil)

	if err != nil {
		log.Printf("RemovePermissionsFromRole failed! %v", err)
		return err
	}

	if resp.UnknownRoleError != nil {
		log.Printf("⚠️ [%s|%v] no such permission exists.", role, perms)
	}
	return nil
}
