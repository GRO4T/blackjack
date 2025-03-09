import { Dispatch, SetStateAction, useEffect } from "react";
import { BASE_URL } from "../constants";
import { GameState, Player } from "../App";

interface Props {
  onGameStartedChanged: Dispatch<SetStateAction<boolean>>;
  gameId: string;
  playerId: string;
  gameState: GameState;
  gameStateSeq: number;
  onGameStateSeqChanged: Dispatch<SetStateAction<number>>;
}

export default function Lobby({
  onGameStartedChanged,
  gameId,
  playerId,
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

  const Leave = async () => {
    await fetch(BASE_URL + "/tables/players/" + gameId + "/" + playerId, {
      method: "DELETE",
    });
    onGameStartedChanged(false);
  };

  return (
    <>
      <div className="column">
        <div id="table-title" className="row centered light-border mid-font">
          Table No. {gameId}
        </div>
        <div id="player-grid" className="row centered">
          {gameState.players &&
            gameState.players.map((player: Player) => (
              <div
                id="player-card"
                key={player.name}
                className="column centered light-border small-font"
              >
                <p>{player.name}</p>
                <input
                  id="player-card-checkbox"
                  type="checkbox"
                  checked={player.isReady}
                  readOnly
                />
              </div>
            ))}
        </div>
        <div className="row centered">
          <div className="column">
            <button onClick={ReportReadiness}>Ready</button>
            <button onClick={Leave}>Leave</button>
          </div>
        </div>
      </div>
    </>
  );
}
