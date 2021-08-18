
/*��� ���������� �������� ������� calculator ���������� ����:
func calculator(arguments <-chan int, done <-chan struct{}) <-chan int
� �������� ��������� ��� ������� �������� ��� ������ ������ ��� ������, ���������� ����� ������ ��� ������.
����� ����� arguments ������� ������� ��� �����, � ����� ����� done - ������ � ������������� ��������� ������.
����� ������ � ���������� ������ ����� �������, ������� ������ � �������� (������������) ����� ���������
����� ���������� �����.
������� calculator ������ ���� �������������, ����� ��������� ����������.
�������� ����� ������ ���� ������ ����� ���������� ���� ����������� �������,
���� �� ����� �� ��������, �� ��������� ���������� ����� ������.*/

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