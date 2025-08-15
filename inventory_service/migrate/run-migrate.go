package migrate

import (
	"fmt"
	"os"
	"os/exec"
)

func RunMigrations(connString string) error {
	// Формируем команду для выполнения миграций
	path := "D:/code/first_krosovochniy_site/backend/inventory_service/migrate/migrations"
	cmd := exec.Command(
		"migrate",           // Исполняемый файл migrate (должен быть в PATH)
		"-path", path, // Путь к папке с миграциями
		"-database", connString, // Строка подключения к БД
		"up",                   // Направление миграции (вверх)
	)

	// Перенаправляем вывод команды в текущий процесс
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	// Запускаем команду и ждём завершения
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ошибка при выполнении миграций: %v", err)
	}

	return nil
}