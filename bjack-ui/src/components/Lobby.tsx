interface Props {
  gameId: string;
}

export default function Lobby({ gameId }: Props) {
  return (
    <>
      Game ID: {gameId}
    </>
  )

}