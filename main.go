package main

import (
	"fmt"
	"strconv"
	"time"
)

func Proceso(id uint64, c chan string, quit chan bool) {
	i := uint64(0)
	for {
		select {
		case <-quit:
			fmt.Println("finalizar")
			return
		default:
			//fmt.Printf("id %d: %d", id, i)
			c <- strconv.FormatUint(id, 10) + ": " + strconv.FormatUint(i, 10)
			i = i + 1
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func SetMostrar(c chan string, notShow chan bool) {
	for {
		select {
		case <-notShow:

			<-c
		case msg := <-c:
			fmt.Println(msg)
		}
	}
}

func SetInactive(notShow chan bool, setActive chan bool) {
	for {
		select {
		case <-setActive:
			return
		default:
			notShow <- true

		}
	}
}

func main() {
	processCount := 0
	chansQuit := make([]chan bool, 1)
	var idCount uint64 = 0
	var opt uint64 = 1
	c := make(chan string)
	notShow := make(chan bool)
	setActive := make(chan bool)
	show := false

	go SetInactive(notShow, setActive)
	for opt != 0 {
		fmt.Println("1.- Agregar proceso")
		fmt.Println("2.- Mostrar procesos")
		fmt.Println("3.- Eliminar proceso")
		fmt.Println("0.- Salir")
		fmt.Scanln(&opt)

		switch opt {
		case 1:
			processCount++
			idCount++
			if idCount > 1 {
				chansQuitTemp := append(chansQuit, make(chan bool))
				chansQuit = chansQuitTemp
			} else {
				chansQuit[0] = make(chan bool)
			}
			go Proceso(idCount, c, chansQuit[idCount-1])
			if idCount == 1 {
				go SetMostrar(c, notShow)
			}
			fmt.Println("Proceso creado")

		case 2:
			if processCount > 0 {
				show = !show
				if show {
					setActive <- true
				} else {
					go SetInactive(notShow, setActive)
				}
			} else {
				fmt.Println("No existen procesos")
			}

		case 3:
			if processCount > 0 {
				var id int64
				fmt.Scanln(&id)
				chansQuit[id-1] <- true
				processCount--
			} else {
				fmt.Println("No existen procesos")
			}
		}
	}
}
