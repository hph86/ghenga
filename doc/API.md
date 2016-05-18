# Introduction

This document provides a high-level overview of the REST API. The objects
ghenga manages are described in the file [Models.md](Models.md).

All objects have a `version` attribute, which is automatically incremented with
each new version of the same object. On update (via `PUT`), the latest version
must be submitted. If the version in the database has increased in the
meantime, the update fails and the user can be informed that someone else has
modified the same object.

When an object is updated, all fields of the object must be submitted in the
`PUT` request. Fields are not present in the JSON data are deleted or reset to
their default value.

# Endpoints

The API is reachable at the path `/api`.

## Authentication

All requests to the API (except the next one) must be authenticated.

### GET /login/token

Log into ghenga with the given user name and password in the HTTP basic auth.
Returns an authentication token which is valid for the given period of time.
The body of response looks as follows:

```json
{
  "user": "foobar",
  "token": "8890bb0467cfe0bde7ec8554b6b01e4174ee6217ed540fc811ef4bfac80c082e",
  "valid_for": 7200,
  "admin": false,
}
```

The token needs to be submitted in the HTTP header `X-Auth-Token` for all
requests to the API.

If the login was not successful, the HTTP response code is 401 (Unauthorized)
and the body will contain a JSON error document.

### GET /login/info

This endpoint can be called with a valid authentication token in the HTTP
header. If the token is still valid, information about the current user and the
remaining validity period is returned. The JSON body is the same as with
`/login/token` endpoint.

### GET /login/invalidate

Performing a GET request to this endpoint invalidates the session token sent in
the `X-Auth-Token` HTTP header. A response code of 200 (OK) and an empty body
is returned on success.

## People

This endpoint manages all entries for people in the database. People can be
communicated with and are assigned to a company.

### GET /person

Returns a list of all persons.

### POST /person

Create a new person. In the body, a JSON document describing the new person
must be submitted. The server responds with a status code of 201 (Created) and
a JSON document with all the data for the new person record, including the ID.

### GET /person/:id:

Returns the data for the specified person.

### PUT /person/:id:

Updates the entry for the person with the specified ID. The body must contain a
JSON document with the changed attributes. Attributes that are not specified
here will not be modified.

### DELETE /person/:id:

Removes the person with the given ID from the database.

## Search

Searching within the data stored by ghenga can be achieved with the following
API endpoints.

### GET /search/person?query=X

This endpoint searches within all people in the database for the string `X`,
where the string may be contained in any field. The response is an array of all
matching people.

## Users

This endpoint manages ghenga users. All requests require the `admin` flag in
the database to be set.

### GET /user

Returns a list of all users.

### POST /user

Create a new user. In the body, a JSON document describing the user
must be submitted.

### GET /user/:id:

Returns the data for the specified user.

### PUT /user/:id:

Updates the entry for the user with the specified ID. The body must contain a
JSON document with the changed attributes. Attributes that are not specified
here will not be modified.

# Errors

When an error occurs, the server returns an appropriate HTTP response code and
an optional JSON document in the body.

For example when ghenga is unable to reach the database server, the HTTP status
code 500 (internal server error) and the following document is returned:

```json
{
  "message": "Unable to connect to database",
  "code": "DATABASE_DOWN"
}
```
