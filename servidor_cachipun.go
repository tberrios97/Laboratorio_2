package main

import (
  "fmt"
  "net"
  "strings"
  "math/rand"
  "strconv"
  "time"
)

//logica aleatoria de jugada
//probabilidad de no estar disponible

func random(min, max int) int {
  //49152-65535
  return rand.Intn(max-min) + min
}

func main() {
  rand.Seed(time.Now().UnixNano())
  fmt.Println("[*] Bienvenidos al servidor de Super-Cachipun 2021")

  PUERTO := ":50000"
  BUFFER := 1024

  estado, error := net.ResolveUDPAddr("udp4", PUERTO)
  if error!=nil {
    fmt.Println(error)
    return
  }

  conexion, error := net.ListenUDP("udp4", estado)
  if error != nil {
    fmt.Println(error)
    return
  }

  defer conexion.Close()
  buffer := make([]byte, BUFFER)

  for {
    n, direccion, error := conexion.ReadFromUDP(buffer)
    //fmt.Println("->", string(buffer[0:n]))

    if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
      //debo enviar OK a servidor intermedio
      fmt.Println("[*] El cliente quiere cerrar la conexion")
      fmt.Println("[*] Se enviará: OK")
      fmt.Println("[*] Apagando servidor")
      mensaje := []byte("OK")
      _, error = conexion.WriteToUDP(mensaje, direccion)
      if error != nil {
        fmt.Println(error)
        return
      }
      return
    }

    if strings.TrimSpace(string(buffer[0:n])) == "DISPONIBLE ?" {
      prob_juego := random(1,10)
      if prob_juego == 1 {
        mensaje := []byte("NO")
        _, error = conexion.WriteToUDP(mensaje, direccion)
        if error != nil {
          fmt.Println(error)
          return
        }
      }else {
        numero_puerto := strconv.Itoa(random(50500, 65535))
        PUERTO_ALEATORIO := ":" + numero_puerto
        estado_c_aleatoria, error := net.ResolveUDPAddr("udp4", PUERTO_ALEATORIO)
        //deberia agregar un try and catch
        if error!=nil {
          fmt.Println(error)
          return
        }

        conexion_aleatoria, error := net.ListenUDP("udp4", estado_c_aleatoria)
        if error != nil {
          fmt.Println(error)
          for {
            numero_puerto = strconv.Itoa(random(50500, 65535))
            PUERTO_ALEATORIO = ":" + numero_puerto
            estado_c_aleatoria, error = net.ResolveUDPAddr("udp4", PUERTO_ALEATORIO)
            conexion_aleatoria, error = net.ListenUDP("udp4", estado_c_aleatoria)
            if error!=nil {
              fmt.Println(error)
            } else {
              break
            }
          }
        }

        fmt.Println("[*] Se iniciara el ejecutor de partidas en el puerto ", numero_puerto)
        mensaje := []byte("OK-localhost-" + numero_puerto)
        _,error = conexion.WriteToUDP(mensaje, direccion)
        if error != nil {
          fmt.Println(error)
          return
        }

        for {
          n, direccion_aleatoria, error := conexion_aleatoria.ReadFromUDP(buffer)
          fmt.Println("[*] El cliente envio: ", string(buffer[0:n]))

          if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
            break
          }

          if strings.TrimSpace(string(buffer[0:n])) == "JUGADA" {
            mensaje := []byte(strconv.Itoa(random(1,3)))
            fmt.Println("[*] Se enviará: ", string(mensaje))
            _, error = conexion_aleatoria.WriteToUDP(mensaje, direccion_aleatoria)
            if error != nil {
              fmt.Println(error)
              break
            }
          }
        }

        conexion_aleatoria.Close()
      }
    }
  }
}
