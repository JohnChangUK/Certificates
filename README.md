# Certificates Back-End Task

## Implemented Endpoints

```$xslt
Method: GET     /certificates                       "Gets all existing certificates"
Method: POST    /certificates                       "Creates a new certificate"
Method: GET     /certificates/{id}                  "Gets a certificate by certificate ID"
Method: PUT     /certificates/{id}                  "Updates a certificate by certificate ID"
Method: DELETE  /certificates/{id}                  "Deletes a certificate by certificate ID"
Method: POST    /certificates/{id}/transfers        "Creates a transfer request to User B"
Method: PUT     /certificates/{id}/transfers        "User B accepts received transfer request"
Method: PATCH   /certificates/{id}/transfers        "User B declines received transfer request"
```

## Use of API
### Creating a certificate

- When creating a new certificate, make a POST request to the endpoint `localhost:8080/certificates`
and provide the `Title`, `Year` and `Note`.
- You must pass in the headers the `Authorization` value, which is the `User Id`. In this instance, it's `John`
```$xslt
{
	"title": "New Verisart certificate",
	"year": 2019,
	"note": "New note"
}
```

The Backend generates a UUID for the certificate Id, timestamp creation, places the 
User Id as the certificate owner Id and an empty Transfer object
```$xslt
{
    "id": "bhpi8r6db2297ufreulg",
    "title": "New Verisart certificate",
    "created_at": "2019-02-24T23:10:36.090193214Z",
    "owner_id": "John",
    "year": 2019,
    "note": "New note",
    "transfer": {
        "email": "",
        "status": ""
    }
}
```

## Scenarios
- User A creates a certificate

1. User A fills in the certificate details on the front-end and creates a certificate

- User A updates a certificate
1. User A views a list of their certificates, clicks on one to open the edit page
2. The user updates the relevant information and saves the certificate

- User A transfers a certificate
1. User A views a list of their certificates, clicks on one to transfer
2. The user enters the information of the person (User B) transferring to and submits.

- User B accepts transferred certificates
1. User B receives an email with a link to accept the transferred certificate
2. User B opens the accept transfer page and clicks "Accept" to finish the transfer
3. User B views a list of their certificates which includes the newly accepted certificate.