# Certificates Back-End Task

## Implemented Endpoints

```$xslt
Method: GET     /certificates                 "Gets all existing certificates"
Method: POST    /certificates                 "Creates a new certificate"
Method: GET     /certificates/{id}            "Gets a certificate by certificate ID"
Method: PUT     /certificates/{id}            "Updates a certificate by certificate ID"
Method: DELETE  /certificates/{id}            "Deletes a certificate by certificate ID"
Method: POST    /certificates/{id}/transfers  "Creates a transfer request to User B"
Method: PUT     /certificates/{id}/transfers  "User B accepts received transfer request"
Method: PATCH   /certificates/{id}/transfers  "User B declines received transfer request"
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