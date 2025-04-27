# Architorture

What started as a card game created by my sister for a school project turned into a fun and creative
collaboration. Together, we transformed her idea into a fully playable digital web version of the game.
The game combines strategy and fun, and now anyone can experience it online.

Check it out here: https://architorture.app.alistairfink.com

## Build

### Configuration

Before building the env vars in the backend and frontend should be configured.

In `Backend/Constants/Constants.go` you should configure the DB variables to the appropriate values. Note: the `DBConnectionString` depends on the IP set in the `docker-compose` file.
```go
DBName             = "architorture"
DBConnectionString = "172.18.0.20:5432"
DBUser             = "postgres"
DBPass             = "replace_with_password"
```

For the frontend the `BaseUrl` const should be set in the `Frontend/src/js/constants/Constants.js` file. This should be whatever the backend enviornment is running on.
```js
const BaseUrl = "architorture-api.app.alistairfink.com";
```

### Images

The `Dockerfile` in the root of this repository can be used to build this project like this. This will create an image
with both the frontend and backend

In order to build the DB image the `Dockerfile` in the `Backend/DatabaseImporter` directory can be used

To run the project the `docker-compose` file can be used
