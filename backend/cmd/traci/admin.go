package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/mail"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"

	"traci/backend/internal/application"
)

func readAdminInput(input io.Reader, output io.Writer) (application.Register, error) {
	reader := bufio.NewReader(input)
	read := func(label string) (string, error) {
		if _, err := fmt.Fprint(output, label+": "); err != nil {
			return "", err
		}
		value, err := reader.ReadString('\n')
		if err != nil && !(errors.Is(err, io.EOF) && value != "") {
			return "", fmt.Errorf("read %s: %w", label, err)
		}
		return strings.TrimSuffix(strings.TrimSuffix(value, "\n"), "\r"), nil
	}
	var registration application.Register
	for _, field := range []struct {
		label string
		value *string
	}{
		{"Никнейм", &registration.Username}, {"Имя", &registration.FirstName},
		{"Фамилия", &registration.SecondName}, {"Почта", &registration.Email},
	} {
		value, err := read(field.label)
		if err != nil {
			return application.Register{}, err
		}
		*field.value = strings.TrimSpace(value)
		if *field.value == "" {
			return application.Register{}, fmt.Errorf("%s не может быть пустым", field.label)
		}
	}
	if utf8.RuneCountInString(registration.Username) > 50 {
		return application.Register{}, fmt.Errorf("никнейм должен содержать не более 50 символов")
	}
	address, err := mail.ParseAddress(registration.Email)
	if err != nil || address.Address != registration.Email || utf8.RuneCountInString(registration.Email) > 255 {
		return application.Register{}, fmt.Errorf("некорректная почта")
	}
	if file, ok := input.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		if _, err := fmt.Fprint(output, "Пароль: "); err != nil {
			return application.Register{}, err
		}
		password, err := term.ReadPassword(int(file.Fd()))
		fmt.Fprintln(output)
		if err != nil {
			return application.Register{}, err
		}
		registration.Password = string(password)
	} else {
		password, err := read("Пароль")
		if err != nil {
			return application.Register{}, err
		}
		registration.Password = password
	}
	if length := utf8.RuneCountInString(registration.Password); length < 8 || length > 128 {
		return application.Register{}, fmt.Errorf("пароль должен содержать от 8 до 128 символов")
	}
	return registration, nil
}
