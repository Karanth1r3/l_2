Что выведет программа? Объяснить вывод программы.

package main

import (
    "fmt"
)

func main() {
    a := [5]int{76, 77, 78, 79, 80}
    var b []int = a[1:4]
    fmt.Println(b)
}
Ответ:

Создается int-массив длины 5
затем создается срез на основе этого массива, в который входят элементы с 1 по 4 (не включительно), т.е. 77, 78, 79;
срез будет иметь длину 3 и ёмкость 4 (ёмкость считается от первого элемента среза до конца массива)
Вывод: [77, 78, 79]