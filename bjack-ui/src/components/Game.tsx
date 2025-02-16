import { GameState, Player, Card } from "../App";

interface Props {
  gameState: GameState,
  playerName: string
}

export default function Game({ gameState, playerName }: Props) {
  return (
    <>
      <div>
      Dealer
      {gameState.hands && gameState.hands[0].map((card: Card) => (
        <div>{card.rank} {card.suit}</div>
      ))}
      </div>
      {
        gameState.players && gameState.players.map((player: Player, index: number) => (
          <div>
            {player.name}
            {gameState.hands[index + 1].map((card: Card) => (
              <div>{card.rank} {card.suit}</div>
            ))}
          </div>
        ))
      }

      {gameState.players[gameState.currentPlayer - 1].name === playerName && 
        <div>My Turn</div>
      }

      State: {gameState.state} <br />
      CurrentPlayer: {gameState.currentPlayer}
    </>
  );
}
