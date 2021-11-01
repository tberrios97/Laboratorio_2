package main

import (
  "fmt"
  "time"
  //"strings"
  "math/rand"
)

func random(min, max int) int {
  //49152-65535
  return rand.Intn(max-min) + min
}

func juego_etapa_1() bool{

}

func juego_etapa_2() bool{

}

func juego_etapa_3() bool{

}

func main(){
  var input int
  var pozo int
  var jugando bool
  rand.Seed(time.Now().UnixNano())
  
  fmt.Println("[*] Bienvenido a SquidGame.\n[*] ¿Deseas unirte?.\n[*] (1) Si.\t(2) No.")
  fmt.Scan(&input)

  //Enviar petición al Lider si se desea jugar. En otro caso, cerrar programa.
  if input == 1 {
    //Enviar petición
    fmt.Println("[*] Iniciando contacto con el Líder.")
    //Establecer conexión al juego y recibir cuando este listo para jugar

    //...

    //Comenzar primera etapa
    jugando = juego_etapa_1()

    if !jugando{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    intermedio("primera")

    jugando = juego_etapa_2()

    if !jugando{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    intermedio("segunda")

    jugando = juego_etapa_3()

    if !jugando{
      fmt.Println("[*] Finalizando programa de SquidGame.")
      return
    }

    pozo = 100
    fmt.Println("[*] Feliciticaciones, has sido uno de los ganadores.\n[*] Has ganado", pozo, "KRW")

  }

  //Cerrar programa
  fmt.Println("[*] Finalizando programa de SquidGame.")
  return
}