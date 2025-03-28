import { useState, Dispatch, SetStateAction } from "react";
import { API_URL, INITIAL_GAME_STATE } from "../constants";
import { GameState } from "../App";

interface Props {
  gameId: string;
  onGameIdChange: Dispatch<SetStateAction<string>>;
  playerName: string;
  onPlayerNameChange: Dispatch<SetStateAction<string>>;
  onGameStartedChange: Dispatch<SetStateAction<boolean>>;
  onPlayerIdChange: Dispatch<SetStateAction<string>>;
  onGameStateChange: Dispatch<SetStateAction<GameState>>;
}

export default function MainMenu({
  gameId,
  onGameIdChange,
  playerName,
  onPlayerNameChange,
  onGameStartedChange,
  onPlayerIdChange,
  onGameStateChange,
}: Props) {
  const [info, setInfo] = useState("");

  const CallCreateGame = async (playerName: string) => {
    return await fetch(API_URL + "/tables", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ playerName: playerName }),
    });
  };

  const CallAddPlayer = async (tableId: string, playerName: string) => {
    return await fetch(API_URL + "/tables/players/" + tableId, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ playerName: playerName }),
    });
  };

  const StartGame = async () => {
    try {
      const createGameResp = await CallCreateGame(playerName);
      if (!createGameResp.ok) {
        if (createGameResp.status == 400) {
          const errMsg = await createGameResp.text();
          setInfo(errMsg);
        }
        throw new Error("POST /tables returned " + createGameResp.status);
      }
      const createGameBody = await createGameResp.json();
      const addPlayerResp = await CallAddPlayer(
        createGameBody["tableId"],
        playerName,
      );
      const addPlayerBody = await addPlayerResp.json();
      onGameStartedChange(true);
      onGameIdChange(createGameBody["tableId"]);
      onPlayerIdChange(addPlayerBody["playerId"]);
      onGameStateChange(INITIAL_GAME_STATE);
    } catch (error) {
      if (error instanceof Error) {
        console.log("Error starting a new game: " + error.message);
      } else {
        console.log("Error starting a new game: " + String(error));
      }
    }
  };

  const JoinGame = async (gameId: string) => {
    try {
      const addPlayerResp = await CallAddPlayer(gameId, playerName);
      if (!addPlayerResp.ok) {
        if (addPlayerResp.status == 404 || addPlayerResp.status == 400) {
          const errMsg = await addPlayerResp.text();
          setInfo(errMsg);
        }
        throw new Error(
          "POST /tables/players/{tableId} returned " + addPlayerResp.status,
        );
      }
      const addPlayerBody = await addPlayerResp.json();
      onGameStartedChange(true);
      onPlayerIdChange(addPlayerBody["playerId"]);
      onGameStateChange(INITIAL_GAME_STATE);
    } catch (error) {
      if (error instanceof Error) {
        console.log("Error joining the game: " + error.message);
      } else {
        console.log("Error joining the game: " + String(error));
      }
    }
  };

  return (
    <>
      <div className="column">
        <input
          value={playerName}
          placeholder="playerName"
          onChange={(e) => onPlayerNameChange(e.target.value)}
        />
        <input
          value={gameId}
          placeholder="tableId"
          onChange={(e) => onGameIdChange(e.target.value)}
        />
        <button onClick={StartGame}>Host a new game</button>
        <button onClick={() => JoinGame(gameId)}>Join a game</button>
        <p className="info">{info}</p>
      </div>
    </>
  );
}
