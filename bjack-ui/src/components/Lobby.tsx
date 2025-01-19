import { useState, useEffect } from "react";
import { BASE_URL } from "../constants";

interface Props {
  gameId: string;
  playerId: string;
}

interface Player {
  name: string;
  isReady: boolean;
  chips: number;
  bet: number;
  outcome: number;
}

interface Card {
  rank: number;
  suit: number;
}

interface GameState {
  players: Player[];
  hands: Card[][];
  state: number;
  currentPlayer: number;
}

export default function Lobby({ gameId, playerId }: Props) {
  const [gameState, setGameState] = useState<GameState>({
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
  }, []);

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
