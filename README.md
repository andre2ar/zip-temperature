# Zip temperature

Fiber web server.

## Running the server

The live reloading server can be run by the following commands:

``make up``

or

``docker-compose up``

There is only one end-point:

``localhost:8080/api/v1/temperature/:zipcode``

Examples request can be found on api folder.

To be able to properly run the server um must add a valid WEATHER_API_KEY on app.env file. Might be necessary to restart
the server after editing app.env file.

GCP URL:

https://zip-temperature-drt52ooijq-uk.a.run.app