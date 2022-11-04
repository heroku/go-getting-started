
# go-getting-started

A barebones Go app, which can easily be deployed to Heroku.

This application supports the [Getting Started with Go on Heroku](https://devcenter.heroku.com/articles/getting-started-with-go) article - check it out.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) version 1.17 or newer and the [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli) installed.

```sh
$ git clone https://github.com/heroku/go-getting-started.git
$ cd go-getting-started
$ go build -o bin/go-getting-started -v . # or `go build -o bin/go-getting-started.exe -v .` in git bash
github.com/mattn/go-colorable
gopkg.in/bluesuncorp/validator.v5
golang.org/x/net/context
github.com/heroku/x/hmetrics
github.com/gin-gonic/gin/render
github.com/manucorporat/sse
github.com/heroku/x/hmetrics/onload
github.com/gin-gonic/gin/binding
github.com/gin-gonic/gin
github.com/heroku/go-getting-started
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

Using resources for this example app counts towards your usage. [Delete your app](https://devcenter.heroku.com/articles/heroku-cli-commands#heroku-apps-destroy) and [database](https://devcenter.heroku.com/articles/heroku-postgresql#removing-the-add-on) as soon as you are done experimenting to control costs.

By default, apps use Eco dynos if you are subscribed to Eco. Otherwise, it defaults to Basic dynos. The Eco dynos plan is shared across all Eco dynos in your account and is recommended if you plan on deploying many small apps to Heroku. Learn more about our low-cost plans [here](https://blog.heroku.com/new-low-cost-plans).

Eligible students can apply for platform credits through our new [Heroku for GitHub Students program](https://blog.heroku.com/github-student-developer-program).

```sh
$ heroku create
$ git push heroku main
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
