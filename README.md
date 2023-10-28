# MobiAPI
WIP scraping library for Polish e-register [mobiDziennik](https://mobidziennik.pl) made in Go

# Examples
Check out [examples/](https://github.com/dark-steveneq/blob/main/examples/) folder for individual examples, however if you want to see reference implementation, check out my abomination of a desktop app - [mobiNG](https://github.com/dark-steveneq/mobing) (hehe funny name). 

# Development
For developing MobiAPI you'll first of all need access to mobiDziennik, a web browser and a basic understanding of Go and probably how various web technologies work and a proxy for intercepting network trafic. I can recommend [ZAP](https://zaproxy.org) since it's what I've been using but it's cross-platform and open source but you can use anything. `MobiAPI` type already has a function for connecting to a proxy server you can use.