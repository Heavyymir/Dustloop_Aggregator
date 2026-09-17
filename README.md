# Dustloop Aggregator
A tool for aggregating fighting game character data from HTML sources. (Wikis)

This is my first personal project for boot.dev. Recent name change due to challenges accessing Supercombo (Cloudflare Challenges) and Mizuumi (Anubis) from command line scrapers.

Repository can be cloned using:

git clone https://github.com/Heavyymir/Dustloop_Aggregator

Disclaimer: I have not gathered or created the data or information being used in this program. This program aims to provide a tool to aggregate it. All credit should go to the 
community members contributing to the data sources used in this program; namely the Dustloop and Mizuumi Wikis. If you would like to support the work of these community members 
through donations etc, please direct them there.

## Project Aims:
The main aim of this project is to create a command line tool with a local SQLite DB instance to quickly and efficiently gather fighting game character frame data in a single local location. 
Ideally, this program will present a lightweight solution for players that only has to touch HTML or JSON data once, or once per patch, that can be used whilst mid set or when labbing to 
reduce bandwidth usage and loading times for data. The use case is niche, but potentially useful. 

Data is to be collected from HTML (Wikis ie. Dustloop, Mizummi) and stored in a Grid - Row - Cell format inside a local SQLite schema.

## Program Commands: 
    • help: Displays a list of usable commands.
	• exit: Exits the program.
	• select: Allows user to select a wiki and game. Usage: select * | select <wiki> * | select <wiki> <game>.
	• fetch: Used after select, fetches a character page and saves to SQLite DB. Usage: fetch <character> [-v|verbose] [-d|details] [-f|--force].
		• [-v|--verbose] Use to print move frame data tables.
		• [-d|--details] Use to print move frame data tables, description and notes.
		• [-f|--force] Use to force scraping of character page and ignore character id check in local DB.
	• discover: Use after select, discovers character names , and saves to a local Json. Usage discover [-n|--no-json].
		• [-n|--no-json] Use to prevent a local Json save.
	• list: Used after selecting a game, lists characters present in the database.
	• frames: Used after selecting a game, Displays framedata for a saved character. Usage frames <character> [-d|--details].
		• [-d|--details] Use to print move description and notes.
	• set path: Allows a user to set the path for Json file saves from the discover command. Usage: set path <directory>.
	• set db: Allows user to switch to a new directory for the internal SQLite DB file. This creates a fresh DB. Usage set db <directory>.
	• dbinfo: Dispalys the location of the internal SQlite DB.


# Usage Considerations:

##Move Notation:
Wikis currently included in this program use Numpad Notation for move inputs. Numpad notation for a move describes moves as if you were pysically moving through a
9 digit numpad on a keyboard, with 5 at the centre (neutral/no motion). Numpad notation assumes you are on the player 1 side (typically left side) of the screen. 
Examples of numpad notation for understanding (Notation listed as numpad motion + x to represent a button press).

	• 236x is a quarter circle forward.
	• 214x is a quarter circle back.
	• 623x is a forward, down, downforward motion.
	• 421x is a back, down, downback motion.
	• 63214 is a half circle back.
	• 41236 is a half circle forward.
	• 632146 is a half circle back to forward.
	• 360/720/1080 are 1, 2 and 3 full circle inputs respectively.
	• [4] 6 is a back charge move.
	• [2] 8 is a down charge move.
	
There are many input types besides these, and this list is not exhaustive.

# Glossary

Fighting game players tend to use specialised language as shorthand for many different concepts. If you require an explanation or clarification on any terminology seen,
please refer to the Fighting Game Glossary at the following link: 

https://glossary.infil.net/

# Initial Program Startup:
Please be aware that this program will prompt you to enter a path to save the required SQLite DB file for local storage of parsed character page HTML data.
When starting the program again, if the DB file is not found at the configured DB path, the program will prompt again for a DB path, and create a brand new DB file.

# When using Fetch:
If attempting to fetch a character that uses hyphens (Zato-1, Jack-o, I-no for Guilty Gear as examples) please include the hyphen "-" to ensure that the
program will correctly build the URL for scraping/assign the character slug correctly in SQLite.

# Package Dependencies:
## Requires
	• github.com/PuerkitoBio/goquery v1.12.0
	• github.com/chromedp/chromedp v0.16.0
	• github.com/chzyer/readline v1.5.1
	• modernc.org/sqlite v1.57.0
	• github.com/andybalholm/cascadia v1.3.3 // indirect
	• github.com/chromedp/cdproto v0.0.0-20260804232424-e85f50dbfd32 // indirect
	• github.com/chromedp/sysutil v1.1.0 // indirect
	• github.com/dustin/go-humanize v1.0.1 // indirect
	• github.com/go-json-experiment/json v0.0.0-20260820222146-c27c302e5fc3 // indirect
	• github.com/gobwas/httphead v0.1.0 // indirect
	• github.com/gobwas/pool v0.2.1 // indirect
	• github.com/gobwas/ws v1.4.0 // indirect
	• github.com/google/uuid v1.6.0 // indirect
	• github.com/mattn/go-isatty v0.0.24 // indirect
	• github.com/ncruces/go-strftime v1.0.0 // indirect
	• github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	• golang.org/x/net v0.52.0 // indirect
	• golang.org/x/sys v0.47.0 // indirect
	• modernc.org/libc v1.74.4 // indirect
	• modernc.org/mathutil v1.7.1 // indirect
	• modernc.org/memory v1.11.0 // indirect



## Development Stretch Goals:
	• Increasing the number of data sources available to include other series, including Street Fighter. This would require access to either SuperCombo (blocked by cloudflare with the
	current fetch methods) or FAT (Would require communication with FAT devs to confirm access to data).
	• Potential inclusion of a GUI version for easier use, but this may contradict the "lightweight" goal of this program.
	• A batch fetch command to fetch all character data for a game at one time.
	• Parsing functions for game mechanics, starter guides and character combos.
	• An internal glossary page to assist new players in understanding terminology used in wikis and by FGC resources in general.
