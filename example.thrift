namespace go example

exception ExampleException {}

service Example {
  i64 Add(1:i64 num1, 2:i64 num2),
  bool Fail() throws (1: ExampleException excp)
}
