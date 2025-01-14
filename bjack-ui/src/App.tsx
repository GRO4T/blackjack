import { useState } from "react";
import "./App.css";
import MainMenu from "./components/MainMenu";
import Lobby from "./components/Lobby";

export default function App() {
  const [gameStarted, setGameStarted] = useState(false);
  const [gameId, setGameId] = useState("");

  if (gameStarted) {
    return <Lobby gameId={gameId} />
  }
  return (
    <MainMenu
      gameId={gameId}
      onGameIdChange={setGameId}
      onGameStartedChange={setGameStarted}
    />
  );
}
