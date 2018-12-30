# simple-backup
A simple way to backup a directory using the storage server with Amazon S3 compatible API like (Minio, Spaces)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

You need Go >= 1.11.2 

### Installing

First you need clone repo

```
go get https://github.com/alejandrojnm/simple-backup
go build
```
other way is download binary for your system, in the release page

## How to use

Copy binary for `/usr/local/bin/` then create cron

`00 2 * * *   root  simple-backup --endpoint=localhost:9000 --ak=accesskey --sk=secretkey --secureendpoint=(true|false) --bucket=system-backup --bucketlocation=us-east-1 --backudir=/srv/backup/mysql`

You can run ``simple-backup -h`` to see help, simple-backup delete all backup from storage server older than 7 days

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Alejandro JNM** - *Initial work* - [alejandrojnm](https://github.com/alejandrojnm)

## License

This project is licensed under the Apache License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc

``