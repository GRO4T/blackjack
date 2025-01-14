import { useState, Dispatch, SetStateAction } from "react";
import { BASE_URL } from "../constants";

interface Props {
  gameId: string;
  onGameIdChange: Dispatch<SetStateAction<string>>;
  onGameStartedChange: Dispatch<SetStateAction<boolean>>;
}

export default function MainMenu({
  gameId,
  onGameIdChange,
  onGameStartedChange,
}: Props) {
  const [info, setInfo] = useState("");

  const StartGame = () => {
    fetch(BASE_URL + "/tables", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: null,
    })
      .then((res) => res.json())
      .then((json) => {
        console.log(json);
        onGameStartedChange(true);
        onGameIdChange(json["tableId"]);
      })
      .catch((error) => {
        console.log("Error starting new game: " + error.message);
      });
  };

  const JoinGame = (gameId: string) => {
    fetch(BASE_URL + "/tables/players/" + gameId, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: null,
    })
      .then((response) => {
        if (!response.ok) {
          if (response.status == 404) {
            setInfo("Game not found");
          }
          throw new Error("POST /tables/players/{tableId} returned " + response.status)
        }
        console.log("Getting JSON");
        response.json();
      })
      .then((json) => {
        console.log(json);
        onGameStartedChange(true);
      })
      .catch((error) => {
        console.log("Error joining the game: " + error.message);
      });
  };

  return (
    <>
      <div className="column">
        <button onClick={StartGame}>Host a new game</button>
        <button onClick={() => JoinGame(gameId)}>Join a game</button>
        <input
          value={gameId}
          placeholder="tableId"
          onChange={(e) => onGameIdChange(e.target.value)}
        />
      </div>
      <p>{info}</p>
    </>
  );
}
