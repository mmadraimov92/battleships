# Battleships

A terminal-based multiplayer Battleships game written in Go.
This project was created as a learning exercise to explore network programming and terminal UI development.

## Features

- **Multiplayer gameplay**: Real-time PvP over network
- **Terminal UI**: Clean ASCII-based interface
- **Ship placement**: Interactive ship positioning with rotation
- **Turn-based combat**: Classic Battleships rules
- **Network architecture**: Client-server model

## How to Play

### Game Setup

The game currently supports PvP mode (AI mode is work in progress).
To start a game, you need to run one server instance and one client instance:

```bash
# Start the server
go run main.go -server -addr 0.0.0.0:1337

# Connect as a client (run in a separate terminal)
go run main.go -addr localhost:1337
```

### Gameplay

1. **Ship Placement**: Each player places their fleet on their board
   - Use arrow keys to move ships
   - Press `R` to rotate ships
   - Press `Enter` to place the ship
   - Ships cannot overlap or go out of bounds

2. **Battle Phase**: Once both players have placed all ships, the game begins
   - The game randomly determines who goes first
   - Players take turns selecting coordinates on the "Target board" to attack
   - Hits and misses are displayed on both boards
   - Continue until one player sinks all opponent ships

3. **Victory**: The first player to sink all enemy ships wins

### Ship Types

The standard fleet consists of:

- 1 × Carrier (5 cells)
- 1 × Battleship (4 cells)  
- 1 × Cruiser (3 cells)
- 1 × Submarine (3 cells)
- 1 × Destroyer (2 cells)

## User Interface

```sh
        My board              Target board
  0 1 2 3 4 5 6 7 8 9     0 1 2 3 4 5 6 7 8 9 
A│~ ~ ~ ~ ~ ~ ~ ~ ~ ~   A│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
B│~ ~ ~ 5 ~ ~ ~ ~ ~ ~   B│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
C│~ ~ ~ 5 ~ ~ ~ ~ ~ ~   C│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
D│~ ~ ~ 5 ~ ~ ~ ~ ~ ~   D│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
E│~ ~ ~ 5 ~ ~ ~ ~ ~ ~   E│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
F│~ ~ ~ 5 ~ ~ ~ ~ ~ ~   F│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
G│~ ~ ~ ~ ~ ~ ~ ~ ~ ~   G│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
H│~ ~ ~ ~ ~ ~ ~ ~ ~ ~   H│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
I│~ ~ ~ ~ ~ ~ ~ ~ ~ ~   I│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 
J│~ ~ ~ ~ ~ ~ ~ ~ ~ ~   J│~ ~ ~ ~ ~ ~ ~ ~ ~ ~ 

Place your ships: Carrier (5 cells)

Controls:
Arrow keys: move ship
R: rotate
Enter: place ship
```

### Legend

- `~` : Water/empty cell
- `#` : Your ship
- `X` : Hit ship
- `O` : Miss
- Numbers: Ship length during placement

## TODO

- [ ] AI mode implementation
- [ ] UI improvements
  - [ ] Better visual indicators/icons
  - [ ] Enhanced status messages
  - [ ] Game state animations
  - [ ] Sound effects (optional)
- [ ] Connection restoration logic
