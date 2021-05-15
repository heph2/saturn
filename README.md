# Saturn-dl

An animesaturn.it scraper and downloader

## Getting Started

These instructions will get you a copy of the project up and
running on your local machine for development and testing
puproses.

### Prerequisities

You need a Go version with modules support, so Go 1.14 at least

### Installing

Just run

``` shell

go get git.mrkeebs.eu/saturn

```

### Basic Usage

``` shell

saturn -f [AnimeToSearch] // Get a list of available episodes with
                          // an index associated
                          
## For Downloading an episode

saturn -d [AnimeToSearch] [id] // This will download the selected ep

```


## Built With

* [Goquery] (https://github.com/Puertokito/goquery/) - Html scraper

## License

This project is licensed under the GPL3 Livence - see the [LICENSE.md](LICENSE.md) file for details
