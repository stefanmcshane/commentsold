# commentsold api
______

* Signin with the following (no other username or password will be accepted)

`curl --location --request GET 'http://localhost:3000/api/signin' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "commentsold",
    "password": "supersecurepassword"
}'`

#### From here and below, we assume that you have set the above token as an env variable by the name 'token'

* List inventory items

`curl --location --request GET 'http://localhost:3000/api/inventory' \
--header 'Authorization: $token`

* Get inventory item by id

curl --location --request GET 'http://localhost:3000/api/inventory/200' \
--header 'Authorization: $token'

* Update inventory item by id to new value

curl --location --request PUT 'http://localhost:3000/api/inventory/200' \
--header 'Authorization: $token' \
--header 'Content-Type: application/json' \
--data-raw '{"ID":200,"ProductID":"9866","Quantity":"27","Color":"Steel Blue","Size":"GIANT","PriceCents":2349,"SalePriceCents":2149}
'