
/*Вам необходимо написать функцию calculator следующего вида:
func calculator(firstChan <-chan int, secondChan <-chan int, stopChan <-chan struct{}) <-chan int
Функция получает в качестве аргументов 3 канала, и возвращает канал типа <-chan int.
В случае, если аргумент будет получен из канала firstChan, в выходной (возвращенный)
канал вы должны отправить квадрат аргумента.
В случае, если аргумент будет получен из канала secondChan, в выходной (возвращенный)
канал вы должны отправить результат умножения аргумента на 3.
В случае, если аргумент будет получен из канала stopChan, нужно просто завершить работу функции.
Функция calculator должна быть неблокирующей, сразу возвращая управление.
Ваша функция получит всего одно значение в один из каналов - получили значение, обработали его, завершили работу.
После завершения работы необходимо освободить ресурсы, закрыв выходной канал, если вы этого не сделаете,
то превысите предельное время работы.
*/
package main

import (
	"fmt"
	"time"
)
import "math/rand"

//---------------------------------------------------------------//
func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	firstChan, secondChan:= make(chan int), make(chan int)
	stopChan := make(chan struct{})
	result := calculator(firstChan, secondChan, stopChan)
	randomNum := rand.Int() % 3
	if randomNum == 0 {
		firstChan <- 3
	} else if randomNum == 1 {
		secondChan <- 3
	} else if randomNum == 2 {
		close(stopChan)
	}
	fmt.Println(<-result)
}

//---------------------------------------------------------------//
func calculator(firstChan <-chan int, secondChan <-chan int, stopChan <-chan struct{}) <-chan int {
	channel := make(chan int)
	go func (firstChan <-chan int, secondChan <-chan int, stopChan <-chan struct{}, channel chan<- int) {
		select {
		case firstChanInt := <-firstChan:
			// Ждем, когда проснется первый канал
			channel <- firstChanInt * firstChanInt
			close(channel)
		case secondChanInt := <-secondChan:
			// Ждем, когда проснется второй канал
			channel <- secondChanInt * 3
			close(channel)
		case <-stopChan:
			// Ждем, когда проснется стоп-канал
			close(channel)
		}
	}(firstChan, secondChan, stopChan, channel)
	return channel
}