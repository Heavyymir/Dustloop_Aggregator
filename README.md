# Dustloop Aggregator
A tool for aggregating fighting game character data from HTML and JSON sources. (Wikis)

This is my first personal project for boot.dev. Recent name change due to challenges accessing Supercombo (Cloudflare Challenges) and Mizuumi (Anubis) from command line scrapers.

Repository can be cloned using:

git clone https://github.com/Heavyymir/Dustloop_Aggregator

Disclaimer: I have not gathered or created the data or information being used in this program. This program aims to provide a tool to aggregate it. All credit should go to the community members contributing to the data sources used in this program; namely Dustloop, FAT, Mizuumi and Supercombo.

Project Aims:
The main aim of this project is to create a command line tool with a local SQLite DB instance to quickly and efficiently gather fighting game character frame data in a single local location. Ideally, this program will present a lightweight solution for players that only has to touch HTML or JSON data once, or once per patch, that can be used whilst mid set or when labbing to reduce bandwidth usage and loading times for data. The use case is niche, but potentially useful. 

Data is to be collected from HTML (Wikis ie. Dustloop, Mizummi, Supercombo) or raw JSON (FAT (Frame Advantage Tool)) and stored in a Grid - Row - Cell format inside SQLite, using Wiki or Fat column headers as the basis to organise and match data.

Potential expansions of project scope include:
    - Increasing the number of data sources to cover more titles.
    - Diversifying data stored. For instance, allowing the storage of wiki game mechanics pages for access by new players.
    - Adding a file output functionality per Game or character to allow for data sheets to be created for local viewing.
    - Parsing combo pages present in Wikis to provide users with the inputs required from community gathered

Usage:

Program Commands: 
    Help: Displays a list of usable commands in the program with usage suggestions
    Exit: Gracefully Exits the program.
    Select <Data Source> <game>: Used to select the target wiki and game. For Example select Dustloop GGST (For Guilty Gear Strive).
    Discover: Used after select, this will search for available Characters for a game.
    Fetch <Character>: Used after select, this will fetch a Character page for the selected wiki and game, then print relevant moves and framedata to the console.
    and store the information in the local SQLite DB. 
    Frames <Character>: Used after select, this will print a table to the console displaying the raw framedata for the selected character, as long as they were
    previously stored in the local SQLite DB.
