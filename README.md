# go-mancer

Create go game engine

## T

| import cycle

```go
Go 언어에서 허용되지 않는 패턴으로, 패키지 A가 패키지 B를 임포트하고, 동시에 패키지 B가 패키지 A를 임포트하는 경우 발생한다.

나의 코드에선 handler 패키지 에서 db 패키지를 임포트하고있고, db 패키지에서 handler 패키지를 임포트하려니 문제가 발생한거다.

이 문제를 해결하기 위해선 해결 방안은 여러가지 있지만, 코드 가독성을 중요시 하는 난 타입을 별도의 패키지에 저장하여서 분산처리를 하게 해주었다.

- handler/auth_handler.go <model에서 타입을 받음>
- db/mongodb.go <model에서 타입을 받음>
- model/register.go <타입을 가지고있음.>

즉, 순환 의존성(circular dependency)를 제거해야한다.
```
