package swaggers

/*
  Exmaple:
	securitySchemes:
    OAuth2AuthCodePKCE:
      type: oauth2
      description: OAuth2 Authorization Code Flow with PKCE support
      flows:
        authorizationCode:
          authorizationUrl: /oauth/authorize
          tokenUrl: /oauth/token
          scopes:
            read: Read access
            write: Write access
*/
func (s *Swagger) OAuth2AuthCodePKCE(AuthorizationUrl, TokenUrl string, Scopes map[string]string) *Swagger {
	s.Components.SecuritySchemes["OAuth2AuthCodePKCE"] = SecurityScheme{
		Type: "oauth2",
		Flows: &OAuthFlows{
			AuthorizationCode: &OAuthFlow{
				AuthorizationUrl: AuthorizationUrl,
				TokenUrl:         TokenUrl,
				Scopes:           Scopes,
			},
			Implicit: &OAuthFlow{

				TokenUrl: TokenUrl,
				Scopes:   Scopes,
			},
			ClientCredentials: &OAuthFlow{
				TokenUrl: TokenUrl,
				Scopes:   Scopes,
			},
			Password: &OAuthFlow{
				TokenUrl: TokenUrl,
				Scopes:   Scopes,
			},
		},
	}
	return s

}
