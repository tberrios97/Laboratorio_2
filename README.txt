# Integrantes:

	- Tomás Berríos 201673576-3 P.200
	- Melanie Corletto 201673529-1 P.201
	- Valeria Miranda 201704503-5 P.200

# Consideraciones:

	- Antes de que muera algún jugador el pozo está en 0.
	- Tomamos por ronda como el número de etapa, es decir 1, 2 y 3 para la creación de los archivos.
	- Si todos pierden en la primera etapa no habrá ganador.
	- Se asume que los inputs son ingresados según las instrucciones dadas.
	- Antes de ejecutar el programa es necesario ejecutar el comando: $sudo service rabbitmq-server start
	- Al termino de cada juego, todas las jugadas son eliminadas, al igual que el registro del pozo.
	- Para los jugadores, solo esta disponible ver el pozo al inicio de cada etapa.
	- Asegurarse de que el servicio de rabbitMQ este corriendo en la maquina 59

# Ejecución:

	- Ingrese mediante la consola a cada maquina, a la ruta del repositorio ().
	- En la maquina correspondiente ejecute el siguiente comando (al momento de realizar las pruebas del laboratorio, por cada comando make se utilizaba una consola de la maquina).
		- Maquina 57:
			- make dataNode
			- make lider
		- Maquina 59:
			- make dataNode
			- make pozo
		- Maquina 60:
			- make nameNode
		- Maquina 58:
			- make dataNode
			- make jugadores (cuando se realiza este comando, se debe de tener paciencia dato que demora su ejecución, aprox 1-2 min)
		- Para cerrar el proceso de cualquiera de las instancias, con el comando Ctrl+C puede ser eliminado.
		- Para cerrar el proceso de los bots, es necesario cerrar el proceso del lider.
		- El orden propuesto de las maquinas es mandatorio para una ejecución correcta.