import { Dispatch, SetStateAction, useEffect } from "react";
import { API_URL } from "../constants";
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
  useEffect(() => {
    onGameStateSeqChanged(gameStateSeq + 1);
  }, []); // eslint-disable-line

  const ReportReadiness = async () => {
    return await fetch(API_URL + "/tables/ready/" + gameId + "/" + playerId, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: null,
    });
  };

  const Leave = async () => {
    await fetch(API_URL + "/tables/players/" + gameId + "/" + playerId, {
      method: "DELETE",
    });
    onGameStartedChanged(false);
  };

  return (
    <>
      <div className="column">
        <div className="table-name row centered light-border mid-font">
          Table No. {gameId}
        </div>
        <div className="players row centered">
          {gameState.players &&
            gameState.players.map((player: Player) => (
              <div
                key={player.name}
                className="player column centered light-border small-font"
              >
                <p>{player.name}</p>
                <input
                  className="player-readiness"
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
