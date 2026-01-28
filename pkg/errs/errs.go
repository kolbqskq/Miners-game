package errs

import "errors"

var (
	ErrGameNotFound       = errors.New("Игра не найдена")
	ErrNotEnoughBalance   = errors.New("Недостоточный баланс")
	ErrAlreadyOwn         = errors.New("Уже куплено")
	ErrSessionIsNotActive = errors.New("Истечен срок сессии")
	ErrServer             = errors.New("Ошибка сервера")
	ErrUserNotFound       = errors.New("Пользователь не найден")
	ErrEmailAlreadyExist  = errors.New("Email уже зарегистрирован")
	ErrIncorrectLogin     = errors.New("Неправильный Email или пароль")
	ErrExpireSession      = errors.New("Сессия истекла")
	ErrEmptyRegisterCode  = errors.New("Введите код")
	ErrRegisterCode       = errors.New("Неверный код")
)
