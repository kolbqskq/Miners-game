package errs

import "errors"

var (
	ErrGameNotFound       = errors.New("Игра не найдена")
	ErrNotEnoughBalance   = errors.New("Недостоточный баланс")
	ErrAlreadyOwn         = errors.New("Уже куплено")
	ErrSessionIsNotActive = errors.New("Истечен срок сессии")
	ErrServer             = errors.New("Ошибка сервера")
	ErrUserNotFound       = errors.New("Пользователь не найден")
)
