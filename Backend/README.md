# Architorture-Backend

## Docker File
Note: Before building change the DBPass value in the Constants/Constants.go file.
### Run the following to compile:
```
docker build -t architorture-backend .
```
This will create two images, the build image and the final image. The build image can be removed.

### Create a Docker Network
```
docker network create --driver bridge --subnet=172.18.0.0/16 architorture-network
```

### Start the DB:
```
docker run -d --name architorture-db -e POSTGRES_PASSWORD=replace_with_password -e PGDATA=/var/lib/postgresql/data/pgdata -p 5430:5432 --network architorture-network --ip 172.18.0.20 postgres
```
Replace the password here too. If the --ip flag is changed change it in Constants.go and rebuild.
Note: This is just a raw postgres db and does not contain the tables or data.
In order to get the data create a new database called "Architoture" and use DatabaseImporter/CreateTables.sql to create the tables.
The sql queries to import the data can be generated using the commands:
```
go run DatabaseImporter/main.go CardSeedData.csv 0
go run DatabaseImporter/main.go CardTypeSeedData.csv 1
go run DatabaseImporter/main.go ExpansionSeedData.csv 2
```

### Run the following to run:
```
docker run -d --name architorture-backend -p 5429:8080 --network architorture-network --ip 172.18.0.21 --restart=always architorture-backend
```
