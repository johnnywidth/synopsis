# Synopsis [![CircleCI](https://circleci.com/gh/johnnywidth/synopsis.svg?style=svg)](https://circleci.com/gh/johnnywidth/synopsis) [![codecov](https://codecov.io/gh/johnnywidth/synopsis/branch/master/graph/badge.svg)](https://codecov.io/gh/johnnywidth/synopsis) [![Go Report Card](https://goreportcard.com/badge/github.com/johnnywidth/synopsis)](https://goreportcard.com/report/github.com/johnnywidth/synopsis)

## Composer Package Repository Generator

Synopsis - it is a tool for generate private composer package repository.
It is work like [satis](https://getcomposer.org/doc/articles/handling-private-packages-with-satis.md).
But more faster, because build time not depend from quantity of repository.

## Need install [golang](https://golang.org/doc/install):
    go get github.com/johnnywidth/synopsis

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
