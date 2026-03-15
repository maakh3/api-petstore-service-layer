# Petstore API (Service-Layer)

A Service-Layer Architecture API implementing the OpenAPI specification found here: https://petstore3.swagger.io/.
Will include a repository pattern for communication with a relational database. Both the database and API will be
containerised and spun up using docker-compose.
All database transactions should be ACID with appropriate locking in place depending on the nature of the transaction.
All response code should be integer values in the first iteration.

-------

## Additional requirements
• All code should be written using TDD where appropriate, with the correct approach, top-down vs bottom-up. 
• Logging should be implemented using slog and logged at the right level. 
• Integration/service tests should be implemented to assure flows across the services. 
• Correct response codes should be supplied

## Required knowledge by the end:
• For any given scenario, the ability to provide the correct response code by number and meaning, covering the most common 2XX, 4XX and 5XX codes. 
• Explaining the difference between pessimistic and optimistic locking and why one would be chosen over the other. 
• Explain the pros and cons of Service Layer and Hexagonal Architectures and why one would be more appropriate than the other, based on requirements of a service.
