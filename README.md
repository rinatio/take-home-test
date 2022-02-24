# Run

    $ go mod download
    $ go run api.go

Create a class:

    $ curl --location --request POST '127.0.0.1:8080/classes' \
      --header 'Content-Type: application/json' \
      --data-raw '{
          "name": "foo",
          "start_date": "2022-02-01T00:00:00Z",
          "end_date": "2022-04-01T00:00:00Z",
          "capacity": 10
      }'

Add a booking

    $ curl --location --request POST '127.0.0.1:8080/classes/1/bookings' \
      --header 'Content-Type: application/json' \
      --data-raw '{
          "name": "Test User",
          "date": "2022-02-01T00:00:00Z"
      }'

Read classes

    $ curl --location --request GET '127.0.0.1:8080/classes'

Read class bookings

    $ curl --location --request GET '127.0.0.1:8080/classes/1/bookings'

# Tests

    $ go test

# Known issues

Using full datetime format due to the following issue with JSON binding:
@link https://github.com/gin-gonic/gin/issues/1193
