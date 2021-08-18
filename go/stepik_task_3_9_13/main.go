
/*��� ���������� �������� ������� calculator ���������� ����:
func calculator(firstChan <-chan int, secondChan <-chan int, stopChan <-chan struct{}) <-chan int
������� �������� � �������� ���������� 3 ������, � ���������� ����� ���� <-chan int.
� ������, ���� �������� ����� ������� �� ������ firstChan, � �������� (������������)
����� �� ������ ��������� ������� ���������.
� ������, ���� �������� ����� ������� �� ������ secondChan, � �������� (������������)
����� �� ������ ��������� ��������� ��������� ��������� �� 3.
� ������, ���� �������� ����� ������� �� ������ stopChan, ����� ������ ��������� ������ �������.
������� calculator ������ ���� �������������, ����� ��������� ����������.
���� ������� ������� ����� ���� �������� � ���� �� ������� - �������� ��������, ���������� ���, ��������� ������.
����� ���������� ������ ���������� ���������� �������, ������ �������� �����, ���� �� ����� �� ��������,
�� ��������� ���������� ����� ������.
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
			// ����, ����� ��������� ������ �����
			channel <- firstChanInt * firstChanInt
			close(channel)
		case secondChanInt := <-secondChan:
			// ����, ����� ��������� ������ �����
			channel <- secondChanInt * 3
			close(channel)
		case <-stopChan:
			// ����, ����� ��������� ����-�����
			close(channel)
		}
	}(firstChan, secondChan, stopChan, channel)
	return channel
}