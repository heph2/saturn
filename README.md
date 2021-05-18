# Saturn

An animesaturn.it scraper and downloader.
Actually can download up to 10 episodes concurrently.

## Getting Started

These instructions will get you a copy of the project up and
running on your local machine for development and testing
puproses.

### Prerequisities

You need a Go version with modules support

### Installing

_Just run_


    go get git.sr.ht/~heph/saturn



### Basic Usage

- For search for the availables Anime run


        saturn [-search] <anime>


- For fetching the available Episodes, using the exact string returned by saturn -search


        saturn [-fetch] <anime>


- For downloading a single episode using the ID associated returned by -fetch

        saturn [-fetch] <anime> -down <ID>


- For downloading differents episodes (eg: episode with ID 1 and episode with ID 3)
  Use a comma-separated input

        saturn [-fetch] <anime> -down <ID,ID>


- For downloadind with a range of episodes (eg: from episode with ID 1 to 6)
  Use a dash-separated input

        saturn [-fetch] <anime> -down <ID-ID>


_Example of Usage:_


    $ saturn -search monster


This return:


    Monster
    Monster-Strike
    Monster-Strike-2
    Monster-Strike-3
    Hatsukoi-Monster
    Monster-Girl-Doctor
    Monster-Strike-2018
    Yu-Gi-Oh-Duel-Monsters-ITA
    Monster-Musume-no-Iru-Nichijou
    Monster-Hunter-Stories-Ride-On
    Digimon-Universe-Appli-Monsters
    Monster-Musume-no-Iru-Nichijou-OVA


Now i can use one of this output for fetching the episodes ( i suggest to copy-paste )


    $ saturn -fetch Monster-Strike


This return:

    ID:0 	 https://www.animesaturn.it/ep/Monster-Strike-ep-1
    ID:1 	 https://www.animesaturn.it/ep/Monster-Strike-ep-2
    ID:2 	 https://www.animesaturn.it/ep/Monster-Strike-ep-3
    ID:3 	 https://www.animesaturn.it/ep/Monster-Strike-ep-4
    ID:4 	 https://www.animesaturn.it/ep/Monster-Strike-ep-5
    ID:5 	 https://www.animesaturn.it/ep/Monster-Strike-ep-6
    ID:6 	 https://www.animesaturn.it/ep/Monster-Strike-ep-7
    ID:7 	 https://www.animesaturn.it/ep/Monster-Strike-ep-8


Now let's say that i want to download from episode 2 to 5

    $ saturn -fetch Monster-Strike -down 1-5



This will *concurrently* download the episodes ranging from 1 to 5.


## Built With

* [Goquery] (https://github.com/PuerkitoBio/goquery) - Like jQuery by for Go

## License

This project is licensed under the GPL3 License - see the [LICENSE](LICENSE) file for details
