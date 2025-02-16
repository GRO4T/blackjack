import { Dispatch, SetStateAction } from "react";
import { GameState, Player, Card } from "../App";

interface Props {
  gameId: string;
  playerId: string;
  gameState: GameState,
  onGameStateChanged: Dispatch<SetStateAction<GameState>>;
  gameStateSeq: number;
  onGameStateSeqChanged: Dispatch<SetStateAction<number>>;
}

export default function Game({ gameState }: Props) {
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

      State: {gameState.state} <br />
      CurrentPlayer: {gameState.currentPlayer}
    </>
  );
}
