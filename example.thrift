namespace go example

service Example {
  i64 Add(1:i64 num1, 2:i64 num2),
  i64 TimeoutedAdd(1:i64 num1, 2:i64 num2, 3:i64 timeoutMS),
}
