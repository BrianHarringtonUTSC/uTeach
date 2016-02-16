# uTeach

uTeach is a reddit-like community oriented platform for sharing educational material and resources.

### Requirements
- Golang 1.4+
- GCC

### Installation Instructions
Note: Ensure your GOPATH is correctly setup.

#### As a User
```
go get github.com/UmairIdris/uTeach
```

#### As a Developer
```
# Make the required directories
mkdir -p $GOPATH/src/github.com/UmairIdris
cd $GOPATH/src/github.com/UmairIdris

# Checkout repo from git
git init && git checkout github.com/UmairIdris/uTeach
cd uTeach

# Install the app and run. Add .exe if on windows to filepath. Export $GOPATH/bin to your PATH for convenience.
go install
$GOPATH/bin/uTeach --config_path=sample_config.json
```

### TODO
- SAML login
- HTTPS
- Admin pages + creating new threads
- Comments
- Middleware for logging/recovery (Gorilla, etc)
- Security (CSRF, etc)
- HTTP tests
