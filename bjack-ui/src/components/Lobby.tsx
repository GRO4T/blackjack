import { Dispatch, SetStateAction, useEffect } from "react";
import { BASE_URL, CARDS_DEALT_STATE, FINISHED_STATE } from "../constants";
import { GameState, Player, Card } from "../App";

interface Props {
  onGameStartedChanged: Dispatch<SetStateAction<boolean>>;
  gameId: string;
  playerId: string;
  playerName: string;
  gameState: GameState;
  gameStateSeq: number;
  onGameStateSeqChanged: Dispatch<SetStateAction<number>>;
}

export default function Lobby({
  onGameStartedChanged,
  gameId,
  playerId,
  playerName,
  gameState,
  gameStateSeq,
  onGameStateSeqChanged,
}: Props) {
  const webSocket = new WebSocket(
    "ws://localhost:8080/state-updates/" + gameId
  );

  useEffect(() => {
    onGameStateSeqChanged(gameStateSeq + 1);
  }, []);

  webSocket.onmessage = (event) => {
    if (event.data === "NewState") {
      onGameStateSeqChanged(gameStateSeq + 1);
    }
  };

  const ReportReadiness = async () => {
    return await fetch(BASE_URL + "/tables/ready/" + gameId + "/" + playerId, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: null,
    });
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

  return (
    <>
      <div className="column">
        Game ID: {gameId} <br />
        Players
        <ul>
          {gameState.players &&
            gameState.players.map((player: Player) => (
              <li key={player.name}>
                {player.name}
                <input type="checkbox" checked={player.isReady} readOnly />
              </li>
            ))}
        </ul>
        Hands
        <ul>
          {gameState.hands &&
            gameState.hands.map((hand: Card[], index: number) => (
              <li key={index}>
                {hand.map((card: Card) => (
                  <span key={card.rank + card.suit}>
                    {card.rank} {card.suit}+
                  </span>
                ))}
              </li>
            ))}
        </ul>
        State: {gameState.state} <br />
        CurrentPlayer: {gameState.currentPlayer}
        {gameState.state === FINISHED_STATE ? (
          <div>{GetOutcome(playerName)}</div>
        ) : (
          <button onClick={ReportReadiness}>Ready</button>
        )}
        <button onClick={() => onGameStartedChanged(false)}>Leave</button>
      </div>
    </>
  );
}
