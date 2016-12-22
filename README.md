# Synopsis [![CircleCI](https://circleci.com/gh/johnnywidth/synopsis.svg?style=svg)](https://circleci.com/gh/johnnywidth/synopsis) [![codecov](https://codecov.io/gh/johnnywidth/synopsis/branch/master/graph/badge.svg)](https://codecov.io/gh/johnnywidth/synopsis) [![Go Report Card](https://goreportcard.com/badge/github.com/johnnywidth/synopsis)](https://goreportcard.com/report/github.com/johnnywidth/synopsis)

## Composer Package Repository Generator

Synopsis - it is a tool for generate private composer package repository.
It is work like [satis](https://getcomposer.org/doc/articles/handling-private-packages-with-satis.md).
But more faster, because build time not depend from quantity of repository.

## To run application need install [golang](https://golang.org/doc/install):
    go get github.com/johnnywidth/synopsis
    cd $GOPATH/src/github.com/johnnywidth/synopsis/
    go build
    ./synopsis

## Docker
    docker run
     -p 9091:8080
     --name synopsis
     -v $(pwd)/data:/data
     -v $HOME/.ssh:/root/.ssh
     johnnywidth/synopsis

## Docker ENV Variables
    HOST        ""                  # dy default is empty, you can set `localhost` and use nginx
    FILE        "/data/config.json" # default path for config file with repositories
    THREAD      50                  # default nomber of concurrent processes
    OUTPUT      "/data/output"      # default path for archived packages

## Web access
`http://localhost:9091/` - Info about building packages.
`http://localhost:9091/admin` - Admin panel.

## Example config.json file:
File look like [satis](https://getcomposer.org/doc/articles/handling-private-packages-with-satis.md).
```
{
  "name": "Private Package Repository",
  "homepage": "http://localhost:9091",
  "archive": {
    "format": "zip",
    "directory": "dist",
    "skip-dev": false
  },
  "repositories": [
    {
      "type": "vcs",
      "url": "git@github.com:johnnywidth/synopsis.git"
    }
  ]
}
```

## Supported satis config params:
- archive
  * `directory`
  * `skip-dev`
  * `format` # only `zip`
  * `prefix-url`
- repositories
  * `type` # only `vcs`, `git`, `composer` (but this mean the same)
  * `url`

## Not supported satis config params:
- [Security](https://getcomposer.org/doc/articles/handling-private-packages-with-satis.md#security) and [Authentication](https://getcomposer.org/doc/articles/handling-private-packages-with-satis.md#authentication)
- `require`, `require-all`, `require-dependencies`, `require-dev-dependencies`
- `abandoned`

## License
Licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/johnnywidth/synopsis/blob/master/LICENSE) for the full license text.
