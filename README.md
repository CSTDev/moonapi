# Moon API

An API for the MoonBoard website, allowing access to problems.

In the future it will allow creating of problems and accessing user information.

### Usage

```
	// Create a new session:
	var session = MoonBoard{}
	
	// Login - auth response will be stored as part of the session
	err := session.Login("Username", "Password")
	
	// Create a query using the builder
	builder := query.New()
	query, _ := builder.Filter(query.Benchmarks).MinGrade(query.SixAPlus).Build()
	
	// Use the session and query to return a slice of Problems
	problems, err := moonBoardSession.GetProblems(query)	
```

#### Cli Usage
Build the command line tool using:
```
	go build .\cmd\cli\main.go
```

Then run it using the following, to list the options use the --help flag
```
	.\main.go --help
```

Example:

```
	.\main.go -login -user username -pass password -hs os,a,b -f Benchmarks
```


### Testing
To run unit tests use the following command from the root directory, it will run all tests:
```
	go test ./...
```

### Documentation
Run godoc:
```
	godoc -http=":<port>"
```

Then visit localhost:<port> and you can navigate to the docs for the package