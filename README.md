# Transactions Application #

This is an application that allows creation of accounts and transferring money between accounts.

### Prerequisites ###
You need to have the following installed to run this application:
* Docker
* Docker compose plugin

### How to run ###
You can run this server by running the following command: `docker compose up`

This will run the application on your local host.

### APIs ###

#### Get account by id ####
This returns accounts by id.

```commandline
curl --location --request GET 'http://localhost/accounts/2'
```

Sample response:
Status: 200 OK
```json
{
    "account_id": 2,
    "balance": 2.3
}
```

#### Create account ####
This creates account with a given id and initial balance.

```commandline
curl --location 'http://localhost/accounts' \
--header 'Content-Type: application/json' \
--data '{
    "id": 2,
    "balance": 2.3
}'
```

Sample response:
Status: 204 No Content

#### Create transaction ####
This creates a transaction which transfers amount from one account to another.

```commandline
curl --location 'http://localhost/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "source_account_id": 123,
    "destination_account_id": 456,
    "amount": "100.12345"
}'
```

Sample response:
Status: 204 No Content

### Future Improvements ###
The following are planned improvements:
* Add a lock when reading and updating accounts when creating a transaction so that concurrent requests for same account dont lead to overwritten data
* Add caching for db calls