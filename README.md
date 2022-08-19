# Golang REST server

A simple REST server.

## Endpoints

POST /bank-accounts
```
{
    "account_number": "1111-11"
}
```
POST /bank-accounts/transfer
```
{
    "from": "1111-11",
    "to": "2222-22"
    "amount": 100
}
```
