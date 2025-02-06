# Generic Data Pipeline

Pipeline provides a lightweight, generic way to create and handle data pipelines.

So insted of manually creating and chaining methdos together, you can use the `pipeline` package.

Inspired by the [Alice](https://github.com/justinas/alice/tree/master) package.

Let's go over what it does using some examples.

## Config Builder

Suppose you want a user to define their own config for your application. And the config struct looks something like this:

```go
type Config struct {
	AppName string
	Port    int
	APIKey  string
	Timeout time.Duration
}

func DefaultConfig(appName string, apiKey string) Config {
	return Config{
		AppName: appName,
		Port:    8080,
		APIKey:  apiKey,
		Timeout: time.Second * 10,
	}
}
```

Now instead of the user having to define their own config from scratch each time, you could offer them a way to build out a pipeline.

Essentially, define functions that would take in a Config, and return a Config

```go
func setTimeout(timeout time.Duration) pipeline.Pipe[Config] {
  return func(cfg Config) Config {
    cfg.Timeout = timeout
    return cfg
  }
}

func setPort(port int) pipeline.Pipe[Config] {
  return func(cfg Config) Config {
    cfg.port = port
    return cfg
  }
}
```

Then finally, use the pipeline as a builder to configure the application.

```go
configureApp(
  pipleine.New(
    SetPort(3000),
    SetTimeout(time.Seconds * 60),
  ).Resolve(Defaultconfig("test app", "api-secret"))),
)
```

Of course, you could use the in built method syntax to do the same thing, but the method syntax is limited to the package the config was defined in. Using the package you could build your own pipleine for complex data structures even if you do not have access to create methods on them.


## Middleware Pipleline

You can also use the pipleine to build out the middleware arcitecture. This is also used in the this [router package](https://github.com/BambooRaptor/router).

NOTE: When closing a pipleline, you get two methods: `Resolve` and `Build` they work similarly, but fulfil different purposes.

Use `Resolve` when the data you're handling are data structures that need to be passed from one pipe to the other.

Use `Build` when the pipes themselved take in function that need to embedded within each-other.
