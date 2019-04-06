# cleaner

[![Licence](https://img.shields.io/github/license/krostar/cleaner.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/cleaner.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/cleaner/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/cleaner)
[![Code quality](https://img.shields.io/codacy/grade/fe2948030b304eeb84eeff103edd9b9f/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/cleaner/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/fe2948030b304eeb84eeff103edd9b9f.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/cleaner/dashboard)

A simple yet effective way of handling initialization and cleaning.

## Motivation

I do really think readability is the best quality a program can have. One thing
that often lack of readability and simplicity is the initialization part in the main
where five, ten, or even more components are initialized on hundred of lines and
errors handling is terrible since defer-able compoments needs to be defered.
So I create a package with very few methods:

-   one method to add function to call when the cleaning happen
-   one method to handle the cleaning

For example, let's say we initialize a logger, a database, and an http server.
If we fail to initialize the database and we forget to flush the logger, some logs
may be lost and we may not know why the database failed to init.

## Usage / example

I like to have two files, one name `cmd/main.go` and another called `cmd/init.go`.
In the main file, I put:

```go
func main() {
    // handle the cleaning if something panic
	defer cleaner.Clean(onCleanError)

	var (
		flags             = parseFlags(os.Args[1:])
		cfg               = initConfig(flags.configFile)
		log, statsHandler = initObs(cfg.Observability)
	)

	log.WithField("version", app.Version()).Info("starting app")

    // here I initialize all the components I need to run the server
	var (
		r10kdeployer = initR10KDeployer(cfg.R10KDeployer)
		httpUsecases = initHTTPUsecases(r10kdeployer)
		srv          = initHTTP(cfg.HTTPServer, httpUsecases, log, statsHandler)
	)

	if err := srv.Run(cfg.GracefulShutdownTimeout, syscall.SIGINT, syscall.SIGTERM); err != nil {
		panic(errors.Wrap(err, "unable to start or stop the server"))
	}
}

func onCleanError(err error) {
	fmt.Fprintf(os.Stderr, "a fatal error occured: %v\n", err)
	os.Exit(2)
}
```

In the init.go, I put all the initialization function:

```go
func initObs(cfg obs.Config) (logger.Logger, string, http.HandlerFunc) {
	logger, statsHandler, stopFunc, err := obs.Init(cfg)
	cleaner.Add(stopFunc) // add the function that flush logs, traces, ...
	if err != nil {
		panic(errors.Wrap(err, "unable to initialize obs component"))
	}

	return logger, statsEndpoint, statsHandler
}

func initHTTP(cfg httpapi.Config, usecases httpapi.Usecases, log logger.Logger, statsHandler http.HandlerFunc) *httpapi.HTTP {
	http, err := httpapi.New(cfg, usecases, log,
		httpapi.WithStatsHandler(statsHandler),
	)
	if err != nil {
		panic(errors.Wrap(err, "unable to create http server"))
	}
	return http
}
```

## License

This project is under the MIT licence, please see the LICENCE file.
