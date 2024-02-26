#!/bin/bash

go mod init sql-trainer.com/m;
go mod tidy;
go install;

echo "Все пакеты успешно установлены!"
go run main.go middleware.go;
