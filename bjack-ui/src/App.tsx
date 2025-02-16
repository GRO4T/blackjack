import { useState } from "react";
import "./App.css";
import MainMenu from "./components/MainMenu";
import Lobby from "./components/Lobby";

export default function App() {
  const [gameStarted, setGameStarted] = useState(false);
  const [gameId, setGameId] = useState("");
  const [playerId, setPlayerId] = useState("");
  const [gameStateSeq, setGameStateSeq] = useState(0);

  if (gameStarted) {
    return (
      <Lobby
        gameId={gameId}
        playerId={playerId}
        gameStateSeq={gameStateSeq}
        onGameStateSeqChanged={setGameStateSeq}
      />
    );
  }
  return (
    <MainMenu
      onGameStartedChange={setGameStarted}
      gameId={gameId}
      onGameIdChange={setGameId}
      onPlayerIdChange={setPlayerId}
    />
  );
}
