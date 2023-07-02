# Identity -Reconciliation

# [UNDERSTANDING]- Project Directory 

- appcontext - initialises the resources for the application.
- common - all the common utilities like databases, migrations reside under this folder
- config - config initialisation ( loading from the env files )
    - database/postgres/migrations - all the migration files related to postgres DB reside under this folder.
    - database/postgres/migrations.go - script to run the migrations.
    - database/postgres/postgres.go - init the database 
- models - database model ( schemas )
- persistence - crud operations 
- profiles - env files
- query - raw queries 
- schema - request/response schemas 
- services - layer which is responsible for calling the persistence layer
    - reconciliation_service - handling all the core logic and calling different layers to make the work done
    - app.go - server init
    - router.go - init the routes 
- utils - all the utilities ( example - logger )


# Docker 

- All the dependencies like postgres and everything are already added, and it's just a click away to start and use it.
- command to start the service -  
  - [1]start the DB - `docker-compose up -d database`
  - [2] start the server - `docker-compose up -d server` 
[NOTE] - follow the order of the commands 

# health check CURL 

`curl --location '0.0.0.0:8080/ping' \
--header 'Cookie: redirect_to=%2Fping'`

# identity API CURL 

[NOTE] - all the notes reside the code with proper internal documentation 

`curl --location '0.0.0.0:8080/api/v1/identify' \
--header 'Content-Type: application/json' \
--header 'Cookie: redirect_to=%2Fping' \
--data-raw '{
"email": "vikashvarma9999@gmail.com",
"phone_number": "123456"
}'`


# Resume
- under resume folder 




