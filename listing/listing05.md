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
  В первой строке мейна переменная имеет тип непустого интерфейса, поэтому он никогда не будет nil (т.к. в itable будет информация об интерфейсе),при присваивании этой переменной результату функции тест в поле значения запишется nil, а второе поле останется прежним (не пустым), поэтому эта переменная никогда не пройдет проверку на пустоту