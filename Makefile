jugador:
	go run ./jugador.go
lider:
	go run ./lider.go
nameNode:
	go run ./nameNode.go
dataNode:
	go run ./dataNode.go
pozo:
	go run ./pozo.go
bot:
	go run ./bot.go

maquina1:
	go run ./dataNode.go & disown
	go run ./lider.go & disown
maquina2:
	go run ./dataNode.go & disown
	go run ./pozo.go & disown
maquina3:
	go run ./nameNode.go & disown
maquina4:
	go run ./dataNode.go & disown

jugadores:
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./bot.go &
	go run ./jugador.go
