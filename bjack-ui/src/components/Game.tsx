import { Dispatch, SetStateAction } from "react";
import { GameState, Player, Card } from "../App";
import {
  BASE_URL,
  CARDS_DEALT_STATE,
  FINISHED_STATE,
  SUIT_CLUBS,
  SUIT_DIAMONDS,
  SUIT_HEARTS,
  SUIT_SPADES,
} from "../constants";
import clubs from "../assets/clubs.png";
import diamonds from "../assets/diamonds.png";
import hearts from "../assets/hearts.png";
import spades from "../assets/spades.png";

interface Props {
  gameId: string;
  playerId: string;
  gameState: GameState;
  playerName: string;
  onGameStartedChanged: Dispatch<SetStateAction<boolean>>;
}

export default function Game({
  gameId,
  playerId,
  gameState,
  playerName,
  onGameStartedChanged,
}: Props) {
  const PlayerAction = async (action: string) => {
    return await fetch(
      BASE_URL + "/tables/" + gameId + "/" + playerId + "?action=" + action,
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: null,
      },
    );
  };

  const GetOutcome = (name: string) => {
    const player = gameState.players.find(
      (player: Player) => player.name === name,
    );
    switch (player?.outcome) {
      case 1:
        return "Won";
      case 2:
        return "Lost";
      case 3:
        return "Push";
      default:
        throw new Error(`Unknown outcome: ${player?.outcome}`);
    }
  };

  const Leave = async () => {
    await fetch(BASE_URL + "/tables/players/" + gameId + "/" + playerId, {
      method: "DELETE",
    });
    onGameStartedChanged(false);
  };

  const GetSuitIcon = (suit: number) => {
    const GetImage = (suit: number) => {
      switch (suit) {
        case SUIT_SPADES:
          return spades;
        case SUIT_DIAMONDS:
          return diamonds;
        case SUIT_CLUBS:
          return clubs;
        case SUIT_HEARTS:
          return hearts;
        default:
          throw new Error(`Unknown suit: ${suit}`);
      }
    };
    return <img src={GetImage(suit)} className="suit" />;
  };

  const GetRankSymbol = (rank: number) => {
    if (rank == 1) {
      return "A";
    }
    if (rank <= 10) {
      return rank;
    }
    switch (rank) {
      case 11:
        return "J";
      case 12:
        return "Q";
      case 13:
        return "K";
      case 14:
        return "JOKER";
      default:
        throw new Error(`Unknown rank: ${rank}`);
    }
  };

  return (
    <>
      <div className="table-name row centered light-border mid-font">
        Table No. {gameId}
      </div>
      <div className="row centered">
        <div className="dealer column centered small-font">
          Dealer
          <div className="hand">
            {gameState.hands &&
              gameState.hands[0].map((card: Card, index: number) => (
                <div
                  key={`${card.rank}-${card.suit}-${index}`}
                  className="card light-border"
                >
                  <div className="row">{GetRankSymbol(card.rank)}</div>
                  <div className="suit-outer row centered">
                    {GetSuitIcon(card.suit)}
                  </div>
                  <div className="rank-flipped row">
                    {GetRankSymbol(card.rank)}
                  </div>
                </div>
              ))}
          </div>
        </div>
      </div>
      <div className="players row centered">
        {gameState.players &&
          gameState.players.map((player: Player, index: number) => (
            <div key={player.name} className="column small-font centered">
              <div className="row centered">
                {player.name}
                {gameState.state === FINISHED_STATE && (
                  <div>&nbsp;({GetOutcome(player.name)})</div>
                )}
              </div>
              <div className="hand">
                {gameState.hands[index + 1].map((card: Card) => (
                  <div
                    key={`${card.rank}-${card.suit}-${index}`}
                    className="card light-border"
                  >
                    <div className="row">{GetRankSymbol(card.rank)}</div>
                    <div className="suit-outer row centered">
                      {GetSuitIcon(card.suit)}
                    </div>
                    <div className="rank-flipped row">
                      {GetRankSymbol(card.rank)}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
      </div>
      <div className="row centered">
        {gameState.state === CARDS_DEALT_STATE &&
          gameState.players[gameState.currentPlayer - 1].name ===
            playerName && (
            <>
              <button onClick={() => PlayerAction("hit")}>Hit</button>
              <button onClick={() => PlayerAction("stand")}>Stand</button>
            </>
          )}
        {gameState.state === FINISHED_STATE && (
          <button onClick={Leave}>Leave</button>
        )}
      </div>
      <footer>
        <small>
          <a
            href="https://www.flaticon.com/free-icons/poker"
            title="poker icons"
          >
            Poker icons created by Freepik - Flaticon
          </a>
        </small>
      </footer>
    </>
  );
}
