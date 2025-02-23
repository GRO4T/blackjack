import { useState, useEffect } from "react";
import "./App.css";
import Game from "./components/Game";
import Lobby from "./components/Lobby";
import MainMenu from "./components/MainMenu";
import { CARDS_DEALT_STATE, BASE_URL } from "./constants";
import { useSessionStorage } from "./useSessionStorage";

export interface Player {
  name: string;
  isReady: boolean;
  chips: number;
  bet: number;
  outcome: number;
}

export interface Card {
  rank: number;
  suit: number;
}

export interface GameState {
  players: Player[];
  hands: Card[][];
  state: number;
  currentPlayer: number;
}

export default function App() {
  const [gameStarted, setGameStarted] = useSessionStorage("gameStarted", false);
  const [gameId, setGameId] = useSessionStorage("gameId", "");
  const [playerName, setPlayerName] = useSessionStorage("playerName", "");
  const [playerId, setPlayerId] = useSessionStorage("playerId", "");
  const [gameStateSeq, setGameStateSeq] = useSessionStorage("gameStateSeq", 0);
  const [gameState, setGameState] = useSessionStorage("gameState", {
    players: [],
    hands: [],
    state: 0,
    currentPlayer: 0,
  });

  useEffect(() => {
    fetch(BASE_URL + "/tables/" + gameId)
      .then((res) => res.json())
      .then((body) => {
        setGameState(body);
      });
  }, [gameId, gameStateSeq]);

  if (gameStarted) {
    if (gameState.state === CARDS_DEALT_STATE) {
      return (
        <Game
          gameId={gameId}
          playerId={playerId}
          gameState={gameState}
          playerName={playerName}
        />
      );
    }
    return (
      <Lobby
        onGameStartedChanged={setGameStarted}
        gameId={gameId}
        playerId={playerId}
        playerName={playerName}
        gameState={gameState}
        gameStateSeq={gameStateSeq}
        onGameStateSeqChanged={setGameStateSeq}
      />
    );
  }
  return (
    <MainMenu
      onGameStartedChange={setGameStarted}
      gameId={gameId}
      playerName={playerName}
      onPlayerNameChange={setPlayerName}
      onGameIdChange={setGameId}
      onPlayerIdChange={setPlayerId}
    />
  );
}
