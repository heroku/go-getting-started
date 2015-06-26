
# go-getting-started

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

A barebones Go app, which can easily be deployed to Heroku.  

This application supports the [Getting Started with Go on Heroku](https://devcenter.heroku.com/articles/getting-started-with-go) article - check it out.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ go get -u github.com/heroku/go-getting-started/...
$ cd $GOPATH/src/github.com/heroku/go-getting-started
$ foreman start web
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

You should also install [Godep](https://github.com/tools/godep) if you are going to add any dependencies to the sample app.

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
