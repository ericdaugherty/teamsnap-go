# teamsnap-go #

teamsnap-go is a simple wrapper for the TeamSnap.com API http://developer.teamsnap.com/documentation/apiv3/

teamsnap-go currently only provides support for reading data from the API, although future updates should add the ability to modify data using the API.

## Testing ##

Simple test cases are provided in teamsnap_test.go  Because the API uses authentication, you must provide an Authentication Token from an OAuth2 Handshake to run the tests.  More information here: http://developer.teamsnap.com/documentation/apiv3/authorization/  Once you have the token, you need to set the environment variable.  On OSX:

```
$ AuthToken="ABC123"
$ export AuthToken
$ go test
PASS
ok  	github.com/ericdaugherty/teamsnap-go	0.747s
```

## Usage ##

The TeamSnap API uses Collection+JSON.  So the API interface is stateful.  Here is a simple example of initializing the API and querying the 'me' object:

```
teamSnap := &TeamSnap{AuthToken: authToken}
teamSnap.Initialize()
meResult, err := teamSnap.FetchRoot("me")
```

From here, you can read the data from the 'me' object in the meResult, or query further.  For example to find the active teams for this user:

```
teamsResult, err := teamSnap.Fetch("active_teams", meResult.Collection.Items[0].Links)
```
