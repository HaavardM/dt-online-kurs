# Golang crash-course for Online NTNU

# Setup
To download service-account credentials for the Disruptive Technologies API, download the `key.json` file from
[this link](https://yopass.disruptive-technologies.com/#/f/7cf0606b-b072-4582-82c8-f3beb408d9a7) and place it in the project folder alongside the `go.mod` file. The file must be named `keys.json`. The download is password protected and the password is shared during the course. 

# How to run the code
All code should be run from the project directory (same directory as the `go.mod` file). 

To run the main application:
```
go run main.go
```

# If you want access to DT studio:
Use the script below to invite yourself to the DT Studio project. You will receive an email to create a new DT Studio user. 
```
go run scripts/invite.go my-email@provider.com
```
[Link to Studio Project](https://studio.disruptive-technologies.com/projects/cs6bd32lh064fjopjmeg/)
