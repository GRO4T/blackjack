import { Dispatch, SetStateAction, useEffect } from "react";
import { BASE_URL } from "../constants";
import { GameState, Player, Card } from "../App";

interface Props {
  gameId: string;
  playerId: string;
  gameState: GameState,
  gameStateSeq: number;
  onGameStateSeqChanged: Dispatch<SetStateAction<number>>;
}

export default function Lobby({
  gameId,
  playerId,
  gameState,
  gameStateSeq,
  onGameStateSeqChanged,
}: Props) {
  const webSocket = new WebSocket("ws://localhost:8080/state-updates/" + gameId);

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
                    {card.rank} {card.suit}
                  </span>
                ))}
              </li>
            ))}
        </ul>
        State: {gameState.state} <br />
        CurrentPlayer: {gameState.currentPlayer}
        <button onClick={ReportReadiness}>Ready</button>
      </div>
    </>
  );
}
