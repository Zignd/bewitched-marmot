# bewitched-marmot

:squirrel: Manga sites scraper as a web API

## Sites support

### Already supported

* Manga Here

### Yet to be supported

* Manga Fox
* Manga Reader
* Manga Panda
* Kiss Manga

## Authentication and Authorization

The API uses 

For the Authorization header you need to provided a [basic access authentication](https://en.wikipedia.org/wiki/Basic_access_authentication) in which the username is a registered client id (client.id on the database) and the password is the client secret (client.secret on the database).

### Requesting the first Access Token

    $ curl -v -X POST -H "Authorization: Basic MTIzNDo1Njc4" -d "grant_type=password" -d "username=zignd" -d "password=qwerty09876" localhost:8080/oauth/token
    > {"access_token":"hznffKNxTxe2VQcIylPvYg","expires_in":3600,"refresh_token":"-kKJECUaRhyGrFsKkozQow","token_type":"Bearer"}

### Requesting an Access Token using a Refresh Token

    $ curl -v -X POST -H "Authorization: Basic MTIzNDo1Njc4" -d "grant_type=refresh_token" -d "refresh_token=-kKJECUaRhyGrFsKkozQow" localhost:8080/oauth/token
    > {"access_token":"ieGqd6QcR5qRGGV3YhYRIA","expires_in":3600,"refresh_token":"1GT4o-FhTl6ONpa2vGsCAg","token_type":"Bearer"}

One interesting thing about our OAuth 2.0 server is that the refresh token can only be used once, attempts to use it a second time will result in the following error message:
    
    {"error":"invalid_grant","error_description":"The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client."}

:D