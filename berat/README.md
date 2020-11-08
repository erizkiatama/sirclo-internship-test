# BE - Berat #

This is the solution or implementation for the SIRCLO internship question. 

A simple CRUD Program of weight management. The features of this program covers:
- Add or Edit Weight Data
- See Detail of a Weight Data
- See all of Weight Data

I created this using Go Programming Language with many tools like GorillaMux, Testify, etc. I am intended of using clean architecture for this program but I think it was too overkill. So, I decided to use MVC instead with package models containing all about models including repository and its mocks, package controller containing all about handler and routers, and views containing all the html templates.

## How To Run - Locally ##

Before run this program locally on your computer, please refer to the .env file and change it into this:
```
DB_HOST=127.0.0.1
#DB_HOST=database
DB_DRIVER=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sirclo
DB_PORT=5432
```

Then to run simply enter this command from terminal and open localhost:8080 from your browser.

```
> go mod download
> go run main.go
```

## How To Run - Docker ##

I have dockerized this program using docker and docker-compose to make it easier to build and run. Before run this program locally on your computer, please refer to the .env file and change it into this:
```
#DB_HOST=127.0.0.1
DB_HOST=database
DB_DRIVER=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=sirclo
DB_PORT=5432
```

To run, you only need to enter this from terminal and open localhost:8080 from your browser.
```
> docker-compose up --build // wait until finished
```

## Testing ##

I already made unit test for all the packages available (excluding main & mocks). To run the unit test please enter this command,
```
go test -cover -covermode=atomic $(go list ./... | grep -v mocks)
```
then you will see the test going and the coverage report.
