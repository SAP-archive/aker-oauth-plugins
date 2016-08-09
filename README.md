# Aker OAuth plugins

This repository contains Aker plugins related to OAuth authentication.


## Aker OAuth Authorization Code Plugin

The Aker OAuth authorization code plugin protects remote resources with the authorization code OAuth flow. This plugin performs the necessary checks to verify that the user has been authorized and has all required scopes.

Following is an example of a valid configuration for this plugin.

```yaml
session:
  authentication_key: "o9pTIkOmETOEfekikEs63X89YFfgXasd"
  encryption_key: "cN2uK5Tl9amDjda2ccapYOJETL4/O1yD"
oauth:
  client_id: "oauth-user"
  client_secret: "oauth-password"
  skip_ssl_validation: true
  authorization_url: https://login.bosh-lite.com/oauth/authorize
  token_url: https://login.bosh-lite.com/oauth/token
  redirect_url: http://localhost:8080/authorization
  required_scopes:
    - messages
  optional_scopes:
    - alerts
    - events
```

The plugin uses Cookies to store the user's session. This is needed so that the component can scale horizontally and support the OAuth flow.

The `session.authentication_key` and `session.encryption_key` specify the authentication key and the encryption keys respectively to be used when securing the cookie that stores the user's session. Both need to be exactly `32` bytes long.

The `oauth.client_id` and `oauth.client_secret` properties are used to authenticate against the OAuth server.

The `oauth.authorization_url` specifies the location where users will be redirected to perform login, as per the OAuth authorization code flow. The `oauth.token_url` specifies the endpoint of the OAuth server where a token can be requested from an OAuth code. The `oauth.skip_ssl_validation` property specifies whether the OAuth server's SSL certificate should be verified.

The `oauth.redirect_url` configures the endpoint where the OAuth server will redirect the user after login. You need to have an `aker-oauth-authorization-code-callback` plugin configured on that endpoint. It is important that the endpoint matches exactly the actual endpoint in terms of URL scheme and host:port.

The `oauth.required_scopes` property specifies a list of scopes that the user will need to have in order to be successfully passed through. The `oauth.optional_scopes` specifies scopes that will be requested from the user but will not be explicitly checked.

This plugin sets the following headers to be used by subsequent plugins:

* `X-Aker-Oauth-Token-Access-Token` - The access token of the user.
* `X-Aker-Oauth-Token-Refresh-Token` - The refresh token of the user.
* `X-Aker-Oauth-Token-Type` - The type of token. (Should be `bearer`)
* `X-Aker-Oauth-Token-Expiry` - The expiration date of the token in `UnixNano` format.
* `X-Aker-Oauth-Info-User-Id` - The user's ID.
* `X-Aker-Oauth-Info-User-Name` - The user's Username.
* `X-Aker-Oauth-Info-User-Scopes` - The scopes that the token has been granted to by the OAuth server.


## Aker OAuth Authorization Code Callback Plugin

The Aker OAuth authorization code callback plugin is an inseparable part of the OAuth Authorization Code plugin. It handles callback requests from the OAuth server to complete the authorization code grant flow.

It has the same configuration as the `aker-oauth-authorization-code` plugin. You should use YAML references in the aker configuration to reduce size and chance of error.

```yaml
session: *session_config
oauth: *oauth_config
```
