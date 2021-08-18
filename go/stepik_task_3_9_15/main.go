
/* Необходимо написать функцию func merge2Channels(fn func(int)
int, in1 <-chan int, in2 <- chan int, out chan<- int, n int).
Описание ее работы:
n раз сделать следующее
прочитать по одному числу из каждого из двух каналов in1 и in2, назовем их x1 и x2.
вычислить f(x1) + f(x2)
записать полученное значение в out
Функция merge2Channels должна быть неблокирующей, сразу возвращая управление.
Функция fn может работать долгое время, ожидая чего-либо или производя вычисления.
Формат ввода:
количество итераций передается через аргумент n.
целые числа подаются через аргументы-каналы in1 и in2.
функция для обработки чисел перед сложением передается через аргумент fn.
Формат вывода:
канал для вывода результатов передается через аргумент out.*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const N = 20 // Константа для определения входящей строки

//-------------------------------------------------------------------------------------------//
func main() {
	in1 := make(chan int, N) // второй аргумент функции merge2Channels
	in2 := make(chan int, N) // третий аргумент функции merge2Channels
	out := make(chan int, N) // четвертый аргумент функции merge2Channels

	start := time.Now()
	merge2Channels(fn, in1, in2, out, N + 1)
	for i := 0; i < N + 1; i++ {
		in1 <- i
		in2 <- i
	}

	orderFail := false
	for i, prev := 0, 0; i < N; i++ {
		fromOutChan := <-out
		if prev >= fromOutChan && i != 0 {
			orderFail = true
		}
		prev = fromOutChan
		fmt.Println(fromOutChan)
	}
	if orderFail {
		fmt.Println("Порядок нарушен")
	}
	duration := time.Since(start)
	if duration.Seconds() > N {
		fmt.Println("Время превышено")
	}
	fmt.Println("Время выполнения: ", duration)
}

//-------------------------------------------------------------------------------------------//
func fn(x int) int {
	// первый аргумент функции merge2Channels
	time.Sleep(time.Duration(rand.Int31n(N)) * time.Second)
	return x * 2
}

//-------------------------------------------------------------------------------------------//
func merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	go func() {
		waitGroup := new(sync.WaitGroup)
		in1ChanMutex := new(sync.Mutex)
		in1Result := make(map[int]int)

		for i := 0; i < n; i++ {
			waitGroup.Add(1)
			in1Value := <-in1
			go func(in1Value int, i int, in1ChanMutex *sync.Mutex, waitGroup *sync.WaitGroup) {
				temp := in1Value
				in1ChanMutex.Lock()
				in1Result[i] = temp
				in1ChanMutex.Unlock()
				waitGroup.Done()
			}(in1Value, i, in1ChanMutex, waitGroup)
		}

		in2ChanMutex := new(sync.Mutex)
		in2Result := make(map[int]int)

		for i := 0; i < n; i++ {
			waitGroup.Add(1)
			in2Value := <-in2
			go func(in2Value int, i int, in2ChanMutex *sync.Mutex, waitGroup *sync.WaitGroup) {
				temp := in2Value
				in2ChanMutex.Lock()
				in2Result[i] = temp
				in2ChanMutex.Unlock()
				waitGroup.Done()
			}(in2Value, i, in2ChanMutex, waitGroup)
		}
		waitGroup.Wait()
		go func() {
			for i := 0; i < n; i++ {
				out <- in1Result[i] + in2Result[i]
			}
			close(out)
		}()
	}()
}