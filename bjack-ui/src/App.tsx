import { useEffect, useRef } from "react";
import "./App.css";
import Game from "./components/Game";
import Lobby from "./components/Lobby";
import MainMenu from "./components/MainMenu";
import { API_URL, INITIAL_GAME_STATE, WAITING_FOR_PLAYERS } from "./constants";
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
  const [gameState, setGameState] = useSessionStorage(
    "gameState",
    INITIAL_GAME_STATE,
  );
  const webSocket = useRef<WebSocket | null>(null);

  useEffect(() => {
    fetch(API_URL + "/tables/" + gameId)
      .then((res) => res.json())
      .then((body) => {
        setGameState(body);
      });
  }, [gameId, gameStateSeq]); // eslint-disable-line

  useEffect(() => {
    if (gameId === "") {
      return;
    }
    webSocket.current = new WebSocket(
      import.meta.env.VITE_API_WS_URL + "/state-updates/" + gameId,
    );
  }, [gameId]);

  if (webSocket.current) {
    webSocket.current.onmessage = (event) => {
      if (event.data === "NewState") {
        setGameStateSeq(gameStateSeq + 1);
      }
    };
  }

  if (gameStarted) {
    if (gameState.state === WAITING_FOR_PLAYERS) {
      return (
        <Lobby
          onGameStartedChanged={setGameStarted}
          gameId={gameId}
          playerId={playerId}
          gameState={gameState}
          gameStateSeq={gameStateSeq}
          onGameStateSeqChanged={setGameStateSeq}
        />
      );
    }
    return (
      <Game
        gameId={gameId}
        playerId={playerId}
        gameState={gameState}
        playerName={playerName}
        onGameStartedChanged={setGameStarted}
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
      onGameStateChange={setGameState}
    />
  );
}
