Что выведет программа? Объяснить вывод программы.

package main

type customError struct {
  msg string
}

func (e *customError) Error() string {
  return e.msg
}

func test() *customError {
  {
    // do something
  }
  return nil
}

func main() {
  var err error
  err = test()
  if err != nil {
    println("error")
    return
  }
  println("ok")
}
Ответ:

...
  err = test() .. - Вся опасность в этой строке. Структуры никогда не бывают nil. Здесь присваивается nil указателю на данные, но сравнение с nil на следующей строке происходит с самой структурой, и результат такого сравнения - всегда true. 
  Поэтому вывод : error