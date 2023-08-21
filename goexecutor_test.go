package goexecutor

import (
    "context"
    "fmt"
    "testing"
)

func TestGoexecutor(t *testing.T) {
    co := New(2)

    for i := 1; i <= 10; i++ {
        ctx := context.WithValue(context.Background(), "id", i)
        co.Work(ctx, func(ctx context.Context) {
            fmt.Println(ctx.Value("id").(int) * 2)
        })

    }
    fmt.Println("等待")
    co.Wait()
    fmt.Println("done")

}
