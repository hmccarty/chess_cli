# GoChess

> This project is under heavy construction! Lots of quick fixes and less than 
ideal solutions being implemented while undergoing testing / research.
So explore at your own risk!

GoChess, as the name might suggest, is a chess CLI written entirely in Golang. I am currently building an internal engine that utilizes deep neural networks to evaluate positions and perform moves. Future work will include integration with the Lichess Go Library.

## Motivations
I wanted a quick CLI to play chess while waiting for build servers at work. For a long time, I had a script that used the Lichess api and that worked
moderately enough. Eventually however, I wanted to explore Deep RL networks, so I revamped this project to include a fully fledged engine.

## Current Status
```
gochess: Includes files for development testing.
│
└─── goengine/: Tracks game positions and handles move generation
│   
└─── docker/: Holds Dockerfile and scripts for portable development
│   
└─── tests/: Contains unit tests for above packages
```

Most development is still happening within the `goengine` package. The progress
of this package is tracked here:

### GoEngine TODO
- ~~Create board representation (bitboard)~~
- ~~Add basic, legal move generation~~
- ~~Add endgame detection~~
- Add support for ~~castling~~, ~~promotions~~ and EP
- Allow for PGN metadata and ~~move parsing~~
- Refactor sliding-piece move generation
- Improve documentation
- Add engine test cases
- Supervised training support for position evaluation
- RL training support for position evaluation