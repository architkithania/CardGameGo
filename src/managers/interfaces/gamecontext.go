package interfaces

type GameContext struct {
	GameId  string
	Players map[*Player]bool
	Host    *Player
}
