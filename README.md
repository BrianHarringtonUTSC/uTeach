# uTeach

uTeach is a reddit-like community oriented platform for sharing educational material and resources.

### Requirements
- Golang 1.4+
- GCC (Sqlite3 dependency)
- Python 3.0+ (if running the database helper)

### Installation Instructions
Notes:
- Ensure your GOPATH is correctly setup
- Export $GOPATH/bin to your PATH for convenience
- Add .exe in front of executables if on Windows

#### As a User
```
go get github.com/umairidris/uTeach
$GOPATH/bin/uTeach --config=$GOPATH/src/github.com/umairidris/uTeach/sample_config.json
```

#### As a Developer
```
# Make the required directories
mkdir -p $GOPATH/src/github.com/umairidris
cd $GOPATH/src/github.com/umairidris

# Checkout repo from git
git init && git checkout github.com/umairidris/uTeach
cd uTeach
go get .

# Install the app and run
go install
$GOPATH/bin/uTeach --config=sample/config.json
```

### TODO
- Add 'tags' to threads
- DB index
- HTTPS
- Admin pages
- Middleware for logging/recovery (Gorilla, etc)
- Don't expose sensitive error info through http.Error calls
- Security (CSRF, etc)
- Front end UI work
- Thread visibility
- Thread pinning
- Fix relative config paths

### FAQ

#### What language does uTeach use? Why?
uTeach is developed in the [Go Programming Language](https://golang.org/).

This may seem like a peculiar choice at first. In fact you may have not even heard of Go before.
Let's begin with understand the environment that uTeach is being developed in.
uTeach will be developed by a few students students at a time working for a period of a few months.
After this, a whole new group of students will take over the project and continue working on it.
The choice of tools should be resilient to several developers of varied skills and experiences working on it for short period of times.
Having a statically typed language greatly solves many of these problems. However, we don't want to add too much overhead to decrease developer productivity. In addition, the tools used cannot be complex, it should be easily picked up i.e. developers should become productive in working on the system in 1-2 weeks max.
It should be familiar (not doing too many new things), cross platform and fast.

Go was a language designed at Google designed by pioneers of the field from the ground up to solve the problems that developers and projects are facing today.


Below are a few links that highlight why Go is a good choice for this project:
- [Used by many established companies to write high performant and stable production systems](https://github.com/golang/go/wiki/GoUsers)
- [I keep seeing mature developers using languages like Go, Rust, Scala, and Erlang; how are those different from using the more common Node/JS, Ruby, PHP, and Python?](https://www.reddit.com/r/webdev/comments/2y3cbf)
- [How we built Uber Engineering's Highest Query Per Second Service Using Go](https://eng.uber.com/go-geofence/)
- [Go at Google: Language Design in the Service of Software Engineering](https://talks.golang.org/2012/splash.article)

This may end up being a mistake in the long run, who knows. Perhaps uTeach will be completely rewritten.
My guess is that that large parts of uTeach will be rewritten to use a client side javascript library thus making the server a lean Go REST api.

### Why doesn't uTeach use a client side framework?
At the time of writing this, the landscape of the web is rapidly changing.
Backbone (taught in C09) is on a decline, Angular is being replaced by Angular 2 which is not backwards compatible.
React is popular now, but might be overshadowed by other libraries once Web Components are standardized (Polymer, etc).
Picking a library now might be outdated in only a few months.
Thus, until the landscape has settled, it will be easier to develop the app as a pure server side, then convert it to a client side app down the line.

### Why SQL? Why Sqlite?
SQL is a proven and time tested technology. It also makes defining relationships between multiple tables clear and efficient.
MongoDB is also an option, however its slow when depending on foreign keys and is more effective for simpler document models.

Sqlite made it easy to get going, thus I used it over others. It should be trivial to switch to another SQL db like postgres.

### What do I need to become productive for this project?

- Start by reading everything on [this page](https://golang.org/doc/). It has guides on how to setup and get running with Go.
- Make sure to do the "Tour of Go" to learn the language quickly (you do not need to do the concurrency sections for this web app).
- Familiarize yourself with Javascript and CSS.
- uTeach is roughly based on the [Go Bootstrap Project](http://go-bootstrap.io/). Some technologies are different, but most of it is the same. More importantly, the structure and layout of the app is closely matched. If you are confused on any part of the project, take a look and see how Go Bootstrap does it for better understanding.
- Check out [Go for Pythonistas](http://s3.amazonaws.com/golangweekly/go_for_pythonistas.pdf) for a guide on Go for Python programmers (if you have finished A08/A48).
- [Comprehensive guide on writing Web apps in Go](https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/)
