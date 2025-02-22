import { GameState, Player, Card } from "../App";
import { BASE_URL } from "../constants";

interface Props {
  gameId: string;
  playerId: string;
  gameState: GameState;
  playerName: string;
}

export default function Game({
  gameId,
  playerId,
  gameState,
  playerName,
}: Props) {
  const PlayerAction = async (action: string) => {
    return await fetch(
      BASE_URL + "/tables/" + gameId + "/" + playerId + "?action=" + action,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: null,
      }
    );
  };

  return (
    <>
      <div>
        Dealer
        {gameState.hands &&
          gameState.hands[0].map((card: Card) => (
            <div>
              {card.rank} {card.suit}
            </div>
          ))}
      </div>
      {gameState.players &&
        gameState.players.map((player: Player, index: number) => (
          <div>
            {player.name}
            {gameState.hands[index + 1].map((card: Card) => (
              <div>
                {card.rank} {card.suit}
              </div>
            ))}
          </div>
        ))}
      {gameState.players[gameState.currentPlayer - 1].name === playerName && (
        <div>
          <button onClick={() => PlayerAction("hit")}>Hit</button>
          <button onClick={() => PlayerAction("stand")}>Stand</button>
        </div>
      )}
      State: {gameState.state} <br />
      CurrentPlayer: {gameState.currentPlayer}
    </>
  );
}
