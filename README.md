# task-scheduler-golang
This is the api developed in golang for scheduling tasks and registering new users. Functionality includes

● Register a new user 

● User login 

● Add a Task 

● Edit a Task 

● Delete a Task 

● Get All Tasks for a user 

● Assign a Internal / External user a task by email address. If the user doesn’t exist send them an email to sign up. Once they signup that note should be assigned to them automatically 

This repo also contains the postman collection to test the api endpoints.

## HOW TO RUN

Run `go run main.go` at root of folder to start the application at http://localhost:8080 by default.

If you have docker installed locally then run the following commands instead at the root of project `docker build . -t task-scheduler` -> `docker run -p 8080:8080 -d task-scheduler`
