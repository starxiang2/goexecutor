package goexecutor

import (
    "fmt"
    "testing"
)

func TestGoexecutor(t *testing.T) {
    co := New(5)

    for i := 1; i <= 10; i++ {
        param := map[string]any{"i": i}

        co.Work(param, func(param map[string]any) {
            fmt.Println(param["i"])
        })

    }

    co.Wait()
}
