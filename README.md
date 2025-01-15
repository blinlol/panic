# PANIC!

Card game for two persons


## Rules

1. The deck is split in half between the players
2. The game is divided into rounds. Each round consists of the following stages:
    1. Some of the cards are laid out on the table in five decks. A maximum of five cards of different denominations can be opened at the same time. 
    The remaining cards are placed face down next to each other.
    2. Players simultaneously reveal one card from their closed deck to the middle of the table. 
    This gives us two open decks in the middle of the table.
    3. You can put cards from a hand with a value different from the values of the open cards by exactly 1 on the open decks. For example, you can put 3 and 5 on the 4, you can put the king and 2 on the ace.
    4. Players can place cards from hand to face, maintain 5 face-up cards of different values on their hand, or combine cards on their hand with the same value.
    5. If each player has cards in his hand and there is no way for any player to make the actions from point 4, then we return to point 2.
    6. If one of the players has run out of cards in his hands, then each player needs to choose one of the open decks as fast as possible. You need to choose a smaller deck faster than your opponent.
    7. The players take the remaining cards in their hand, their closed deck and their open deck.
    8. If everyone has cards, then we play the next round.
3. The player with no remaining cards wins.

## Start

To install and run Panic!, run following commands:

```bash
go install github.com/blinlol/panic@latest
panic
```
