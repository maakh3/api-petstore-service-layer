# Petstore API (Service-Layer)

A Service-Layer Architecture API implementing the OpenAPI specification found here: https://petstore3.swagger.io/.
Will include a repository pattern for communication with a relational database. Both the database and API will be
containerised and spun up using docker-compose.
All database transactions should be ACID with appropriate locking in place depending on the nature of the transaction.
All response code should be integer values in the first iteration.