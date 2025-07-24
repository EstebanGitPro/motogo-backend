#!/bin/bash


echo "⏫ Iniciando contenedor MySQL (mysql-motogo)..."
sudo docker start mysql-motogo


echo "⏳ Esperando a que el contenedor inicie..."
sleep 5


echo "🚀 Ejecutando la aplicación Go..."
go run /home/devban/Documents/Go/motogo-backend/cmd/main.go