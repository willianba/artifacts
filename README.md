# Artifacts

Recently I discovered [Artifacts MMO](https://www.artifactsmmo.com/) and saw it as an opportunity to learn Go. I also found many other libraries and SDKs that already do some wrappers around its API so I decided I could build my own bot to play the game.

This repo is literally this: a CLI/Bot to play artifacts (while I learn Go!). Don't expect the code to be perfect.

# How to use

- Create a .env file and store your artifacts token

```
TOKEN=ey.....
```

- Run `go run main.go` for a list of available commands
    - I suggest to read the `main.go` file to understand it better
    - The `-name` flag is mandatory for all actions (it's possible to use `-n` too)

> You can run `go build` to generate the `artifacts` executable too

## Examples

Here are some examples on how to use the CLI. There are many aliases. You could type the words or get used to using their alisases. Check the code for all of them or tun `go run main.go help` and `go run main.go <action> help`.

### Direct commands

- Move: `go run main.go -n JohnDoe ch m 1 2`
    - Read as "Character, Move, (1,2)" to move the character to the task master
- Fight: `go run main.go -n JohnDoe ch f` to fight on the current map
    - Fight also has a `-r` flag to rest after fighting: `go run main.go -n JohnDoe ch r -r`
- Deposit to bank: `go run main.go -n JohnDoe b d -i copper -q 10`
    - Reads as "Bank, Deposit, 10 copper"
- Looping a specific action: `go run main.go -n JohnDoe -l 10 ch g`
    - Read as "Loop 10, Character, Gather" to run 10 gather commands
    - Loops can be added to any command, but I didn't implement loops for all of them (like moving, crafting, and so on)

### Flows

This is my attempt to automate actions (aka create bots). There are little to no flows at the moment, but as I progress in the game I'll creating more of them. To use flows you would do the same as a direct command, but instead of running only one action, it will run a loop of many actions.

This is very error prone and probably has some edge cases, but it works fine until certain point so I'm ok with them.

- Forge a mineral: `go run main.go flow copper 10` to run the copper flow until you reach 10 copper in the inventory
    - You can swap copper for iron too
