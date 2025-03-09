import { Dispatch, SetStateAction } from "react";
import { GameState, Player, Card } from "../App";
import { BASE_URL, CARDS_DEALT_STATE, FINISHED_STATE } from "../constants";

interface Props {
  gameId: string;
  playerId: string;
  gameState: GameState;
  playerName: string;
  onGameStartedChanged: Dispatch<SetStateAction<boolean>>;
}

export default function Game({
  gameId,
  playerId,
  gameState,
  playerName,
  onGameStartedChanged,
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

  const GetOutcome = (name: string) => {
    const player = gameState.players.find(
      (player: Player) => player.name === name
    );
    switch (player?.outcome) {
      case 1:
        return "Win";
      case 2:
        return "Lose";
      case 3:
        return "Push";
      default:
        return "Unknown";
    }
  };

  const Leave = async () => {
    await fetch(BASE_URL + "/tables/players/" + gameId + "/" + playerId, {
      method: "DELETE",
    });
    onGameStartedChanged(false);
  };

  return (
    <>
      <div id="dealer-table" className="row centered">
        <div id="player-hand" className="row centered light-border small-font">
          Dealer
          {gameState.hands &&
            gameState.hands[0].map((card: Card) => (
              <div>
                {card.rank} {card.suit}
              </div>
            ))}
        </div>
      </div>
      <div id="player-grid" className="row centered">
        {gameState.players &&
          gameState.players.map((player: Player, index: number) => (
            <div
              id="player-hand"
              key={player.name}
              className="column light-border small-font"
            >
              {player.name}
              <div id="card-grid">
                {gameState.hands[index + 1].map((card: Card) => (
                  <div
                    id="card"
                    key={card.rank + card.suit}
                    className="light-border"
                  >
                    {card.rank} {card.suit}
                  </div>
                ))}
              </div>
              {gameState.state === FINISHED_STATE && (
                <div>{GetOutcome(player.name)}</div>
              )}
            </div>
          ))}
      </div>
      <div className="row centered">
        {gameState.state === CARDS_DEALT_STATE &&
          gameState.players[gameState.currentPlayer - 1].name ===
            playerName && (
            <>
              <button onClick={() => PlayerAction("hit")}>Hit</button>
              <button onClick={() => PlayerAction("stand")}>Stand</button>
            </>
          )}
        {gameState.state === FINISHED_STATE && (
          <button onClick={Leave}>Leave</button>
        )}
      </div>
    </>
  );
}
