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

### Errors
| Code | ErrorType | Detail |
| ----------- | ----------- | ----------- |
| 101 | InvalidRequestErr | invalid request parameters|
| 102 | InternalErr | Internal Error, you can try again at another time|