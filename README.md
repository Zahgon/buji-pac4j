# buji-pac4j (Go)

A Go implementation of the **buji-pac4j** security bridge: it pushes a
[pac4j](https://www.pac4j.org/) security context (authenticated user profiles,
roles and permissions) into an [Apache Shiro](https://shiro.apache.org/)-style
security context, so that Shiro's `Subject` reflects the identity established by
pac4j.

This is a faithful port of the original Java library. It preserves the same
observable functionality, security semantics, authentication/authorization
flow, configuration behavior, and session lifecycle — reimplemented in idiomatic
Go with zero Java runtime or Maven dependency.

## Requirements

- Go 1.23 or newer.

## Build & Test

```bash
go build ./...
go vet ./...
go test ./...
```

The test suite mirrors the original library's tests one-for-one (17 tests
covering principal naming, the pac4j→Shiro bridge, session renewal, and
configuration).

## What it does

The library bridges two security models:

- **pac4j** produces a set of authenticated `UserProfile`s (one per client),
  each carrying an id, a client name, attributes, and roles.
- **Shiro** exposes a `Subject` with a primary principal, a principal
  collection, authentication state, roles, and permissions.

The bridge takes the pac4j profiles and logs the Shiro `Subject` in with them,
so that downstream Shiro-based authorization sees the pac4j identity.

## Security behavior

### Authentication

`buji/util.PopulateSubject` is the core entry point. Given the ordered pac4j
profiles it:

1. Flattens them into a profile list.
2. Decides, via the pac4j authorizers, whether the user is *fully
   authenticated* (at least one non-remembered profile) or merely *remembered*.
3. Logs the Shiro `Subject` in with a `Pac4jToken` carrying those profiles.

On a real login the Shiro security manager **renews the session id**
(session-fixation protection). The primary principal of the resulting subject is
the computed **username string**; the principal collection also holds the
`Pac4jPrincipal` object (retrievable via `OneByType`).

### Principal name

`Pac4jPrincipal.GetName()` resolves the username:

- If no principal-name attribute is configured, it returns the profile id.
- Otherwise it returns the named attribute's value (via `String.valueOf`
  semantics), or nothing when the attribute is absent.

A configured attribute name is trimmed; a blank attribute name is treated as
unset.

### Profile refresh keeps the session

When `PopulateSubject` is called again for the **same user** — same number of
profiles, and each profile matching by client name and id, in order, with an
unchanged computed principal name — the bridge refreshes the existing
`Pac4jPrincipal` **in place** instead of logging in again. This keeps the
current session id (no re-login) while updating attributes such as a refreshed
access token. Any change to identity (id), client, or principal name triggers a
fresh login and a new session id.

### Authorization

`Pac4jRealm` derives Shiro authorization info from the profiles:

- **Roles** come from each profile's roles.
- **Permissions** come from the profile attribute named by the
  `SHIRO_PERMISSIONS` key (`$shiro_permissions$`), accepting a list or set of
  strings.

## Configuration

The default configuration resource `buji-pac4j-default.ini` wires the pac4j
realm and subject factory into the security manager. `Pac4jIniEnvironment`
loads this embedded resource as the framework configuration.

## Package layout

- `pac4j/` — the pac4j-side model: `profile` (UserProfile/CommonProfile and
  helpers), `authorization/authorizer` (fully-authenticated / remembered),
  `context` and `context/session` (web context and session store contracts),
  `config`, `client`, `cas/config`, `exception`, `framework/adapter`.
- `shiro/` — the Shiro-side model: `subject` (Subject, SubjectContext,
  PrincipalCollection), `session`, `authc`, `authz`, `realm`, `mgt`
  (security manager, delegating subject), `util` (thread-bound security
  manager and `GetSubject`), `web/mgt`, `web/env`, `io` (serializer).
- `buji/` — the bridge itself: `subject` (`Pac4jPrincipal`,
  `Pac4jSubjectFactory`), `token` (`Pac4jToken`), `realm` (`Pac4jRealm`),
  `context` (`ShiroSessionStore`), `profile` (`ShiroProfileManager`),
  `util` (`PopulateSubject` and refresh logic), `env`
  (`Pac4jIniEnvironment`), `resources` (embedded default configuration).
- `internal/adapter` — the framework adapter that installs the bridge defaults
  (profile manager factory and session store factory) into a pac4j config.

## License

Apache License 2.0.
