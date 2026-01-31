package game

import "miners_game/internal/game/domain"

func PutGameToMemory(s *Service, userID, gameID string, game *domain.GameState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.games[userID+"/"+gameID] = game
}
