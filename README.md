# Text Adventure Game in Go

This repository contains a text-based adventure game implemented in Go. The game was created as part of a homework assignment to practice basic coding skills, working with structures, methods, functions, and maps in Go.

## Description

The game is a simple text-based adventure where the player can navigate through different rooms, interact with objects, and perform various actions. The player starts in the kitchen and can explore other rooms like the corridor, bedroom, and outside. The game reacts to player commands and provides appropriate responses.

## Features

- Navigate through different rooms
- Interact with objects in the rooms
- Perform actions such as "look around", "pick up", "wear", and "apply"
- Implemented using structures and functions to allow easy extension of the game world

## Requirements

- Go programming language (version 1.16 or later)

## Installation and Running

1. Clone the repository or download the source code.
2. Navigate to the directory with the source code.
3. Run the game using the Go command:

    ```sh
    go run main.go
    ```

4. To run the tests:

    ```sh
    go test -v
    ```

## How to Play

The game accepts commands from the player in the following format:

```sh
$command $parameter1 $parameter2
```

## Available Commands

```
* осмотреться — Look around the current room.
* идти [direction] — Move to another room (e.g., "идти коридор").
* взять [item] — Pick up an item (e.g., "взять ключи").
* надеть [item] — Wear an item (e.g., "надеть рюкзак").
* применить [item] [object] — Apply an item to an object (e.g., "применить ключи дверь").
* выход — Exit the game.
```

## Example Commands

```
* осмотреться
* идти коридор
* взять ключи
* надеть рюкзак
* применить ключи дверь
```

## Game World

```
The game world consists of the following rooms:

- Kitchen (кухня): The starting room with some items.
- Corridor (коридор): A connecting room with exits to other rooms.
- Bedroom (комната): A room with more items to pick up.
- Outside (улица): The final room you need to unlock by applying the keys to the door.
```
