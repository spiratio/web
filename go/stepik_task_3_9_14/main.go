
/*Вам необходимо написать функцию calculator следующего вида:
func calculator(arguments <-chan int, done <-chan struct{}) <-chan int
В качестве аргумента эта функция получает два канала только для чтения, возвращает канал только для чтения.
Через канал arguments функция получит ряд чисел, а через канал done - сигнал о необходимости завершить работу.
Когда сигнал о завершении работы будет получен, функция должна в выходной (возвращенный) канал отправить
сумму полученных чисел.
Функция calculator должна быть неблокирующей, сразу возвращая управление.
Выходной канал должен быть закрыт после выполнения всех оговоренных условий,
если вы этого не сделаете, то превысите предельное время работы.*/

package main

import "fmt"

//---------------------------------------------------------------//
func main(){
	arguments := make(chan int)
	done := make(chan struct{})
	result := calculator(arguments, done)
	for i := 0; i < 10; i++ {
		arguments <- 1
	}
	close(done)
	fmt.Println(<-result)
}

//---------------------------------------------------------------//
func calculator(arguments <-chan int, done <-chan struct{}) <-chan int {
	channel := make(chan int)
	go func(arguments <-chan int, done<-chan struct{}, channel chan <- int ) {
		var result int
		for {
			select {
			case temp := <-arguments:
				result = result + temp
			case <-done:
				channel <- result
				close(channel)
				return
			}
		}
	}(arguments, done, channel)
	return channel
}