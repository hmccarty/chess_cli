# GoChess

> This project is under heavy construction! Lots of quick fixes and less than 
ideal solutions being implemented while undergoing testing / research.
So explore at your own risk!

GoChess, as the name might suggest, is a chess CLI written entirely in Golang. I am currently building an internal engine that utilizes deep neural networks to evaluate positions and perform moves. Future work will include integration with the Lichess Go Library.

## Motivations
I wanted a quick CLI to play chess while waiting for build servers at work. For a long time, I had a script that used the Lichess api and that worked
moderately enough. Eventually however, I wanted to explore Deep RL networks, so I revamped this project to include a fully fledged engine.

## TODO
- ~~Create board representation (bitboard)~~
- ~~Add basic, legal move generation~~
- ~~Add endgame detection~~
- Add support for ~~castling~~, ~~promotions~~ and EP
- Create engine test cases
- RL Self-Play Support
