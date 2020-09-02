# challenge-modec

[![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/victorabarros/travel-routes-optimizer/master/LICENSE)

back-end challenge from Modec

## Description

### goal

The assignment involves creating a backend to manage different equipment of an FPSO (Floating Production, Storage and
Offloading). This system will be used for other applications in the organization and we should have APIs with the appropriate
HTTP request methods to be able to reuse them. The data should be stored in the database (you can use in-memory database).

### functionalities

- Registering a vessel. The vessel data input is its code, which can’t be repeated (return the HTTP code appropriate and an error message if the user tries to register a existing code). For instance, a valid input of a vessel is: “code”: “MV102”.
- Registering a new equipment in a vessel. The data inputs of each equipment are name, code, location and status. Each equipment is associated to a given vessel and has a unique code, which can’t be repeated (return the HTTP code appropriate and an error message if the user tries to register a existing code). For each new equipment registered, the equipment status is automatically active. For instance, a valid input of a new equipment related to a vessel “MV102” is:

```json
    {
        "name": "compressor",
        "code": "5310B9D7",
        "location": "Brazil"
    }
```

- Setting an equipment’s status to inactive. The input data should be one or a list of equipment code.
- Returning all active equipment of a vessel

### evaluate

- Best practices on how you design your solution
- Unit tests are mandatory
- Software engineering principles: API design, separation of concerns and modularity

## Development

### Programming languange choice

To improve a better performance and explore the native tools, I choose [Golang](https://golang.org/). It's a language getting more space on market, simple writenning and fast.
To improve the infrastructure and make a easy scalability the development was made with [Docker](https://docs.docker.com/).

### Requirements

- [Docker](https://docs.docker.com/)
- [GNU make](https://www.gnu.org/software/make/)

### Endpoints

| Method |           Route            |         Description       |
|--------|----------------------------|---------------------------|
| GET    | `/healthz`                 | liveness probe            |
| GET    | `/healthy`                 | readness probe            |
| POST   | `/vessel`                  | insert vessel             |
| GET    | `/vessel/{code}`           | list equipments by vessel |
| POST   | `/vessel/{code}/equipment` | insert single equipment   |
| POST   | `/vessel/{code}/equipments`| insert list of equipments |
| DELETE | `/equipemnt/{code}`        | inactive equipemnt        |

More details about payloads are at the [collection](./dev/Challenge-Modec.postman_collection.json)

### How to Run

Write `.env` file based on [.env.example](.env.example) and execute:

```sh
make run
```

### Tests

To see the html coverage [c.out](./dev/c.out), execute:

```sh
make test-coverage
```

To see .log , execute:

```sh
make test-log
```

output [tests-summ.log](./dev/tests-summ.log)

```log
coverage: 100.0% of statements
ok  	github.com/victorabarros/challenge-modec/app/server	0.085s	coverage: 100.0% of statements
coverage: 83.3% of statements
ok  	github.com/victorabarros/challenge-modec/internal/config	0.031s	coverage: 83.3% of statements
```
