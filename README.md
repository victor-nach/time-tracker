# Time tracker app
Time-tracker is a small full stack web app, that can help a freelancer track their time.

User stories:

- As a user, I want to be able to start a time tracking session
- As a user, I want to be able to stop a time tracking session
- As a user, I want to be able to name my time tracking session
- As a user, I want to be able to save my time tracking session when I am done with it
- As a user, I want an overview of my sessions for the day, week and month
- As a user, I want to be able to close my browser and shut down my computer and still have my sessions visible to me when I power it up again.


## How to start 

```shell script
# To generate graphql schemas
# Note: see server/graph/schemas and also ./gqlgen.yml
$ make schema 
```

```shell script
# Use this command to start the service locally

$ make local
```

```shell script
# To run tests

$ make test 
```

```shell script
# build a binary for the api

$ make build 
```

```shell script
# generate mocks

$ make gen-mocks 
```

```shell script
# runs go vet

$ make vet 
```

```shell script
# runs go fmt to format files

$ make fmt 
```

## Deployments

- Backend deployed version - https://trackerr-app.herokuapp.com/
- frontend deployed version -  https://victor-nach.github.io/time-tracker-frontend/
- frontend repo - https://github.com/victor-nach/time-tracker-frontend


<img width="1600" alt="Screenshot 2021-06-16 at 03 44 11" src="https://user-images.githubusercontent.com/46886694/122150439-a6264500-ce55-11eb-85ba-1cc6f56c55c2.png">
<img width="1535" alt="Screenshot 2021-06-16 at 03 24 34" src="https://user-images.githubusercontent.com/46886694/122150817-585e0c80-ce56-11eb-8cc4-8619e59f4a15.png">


## Features
- Signup
- Login
- Refresh tokek
- Get user info
- save session
- view saved sessions
- update session
- delete session

# Tools
- Go
- GraphQL
- MongoDB (datastore)
- Heroku (deployment)

### Internal Error definition
| Code | ErrorType | Detail |
| ----------- | ----------- | ----------- |
| 101 | InvalidRequestErr | invalid request parameters |
| 102 | InternalErr | Internal error|
| 103 | DatabaseErr | database error |
| 104 | InvalidAuthErr | email or passcode invalid |
| 105 | CustomerNotFoundErr | invalid customer id |
| 106 | SessionNotFoundErr | invalid session id |
| 107 | EmailExistsError | Duplicate Email found |

