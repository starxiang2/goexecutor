package goexecutor

import (
    "context"
    "errors"
    "sync"
)

var maxGoroutine chan struct{}
var maxGoroutineNum = 0

func SetGlobalMaxGoroutine(maxNum int) error {
    if maxGoroutineNum > 0 {
        return errors.New("不能重复设置全局协程数量")
    }
    maxGoroutineNum = maxNum
    maxGoroutine = make(chan struct{}, maxGoroutineNum)
    return nil
}

func GetGlobalGoroutineCount() (int, error) {
    if maxGoroutineNum == 0 {
        return 0, errors.New("没有设置全局协程数量")
    }
    return maxGoroutineNum, nil
}

func GetCurrentGlobalGoroutineCount() (int, error) {
    if maxGoroutineNum == 0 {
        return 0, errors.New("没有设置全局协程数量")
    }

    return len(maxGoroutine), nil
}

type goroutineControl struct {
    GoroutineNum chan struct{}
    wg           *sync.WaitGroup
}

func (g *goroutineControl) Work(ctx context.Context, f func(ctx context.Context)) {
    g.Add()
    go func(c *goroutineControl, ctx context.Context, fu func(ctx context.Context)) {
        defer c.Done()
        fu(ctx)
    }(g, ctx, f)
}

func (g *goroutineControl) Add() {
    g.wg.Add(1)
    g.GoroutineNum <- struct{}{}
    if maxGoroutineNum > 0 {
        maxGoroutine <- struct{}{}
    }
}

func (c *goroutineControl) GetCurrentGoroutineCount() int {
    return len(c.GoroutineNum)
}

func (g *goroutineControl) Done() {
    <-g.GoroutineNum
    if maxGoroutineNum > 0 {
        <-maxGoroutine
    }
    g.wg.Done()
}

func (g *goroutineControl) Wait() {
    g.wg.Wait()
}

func New(goroutineNum uint16) *goroutineControl {
    return &goroutineControl{
        GoroutineNum: make(chan struct{}, goroutineNum),
        wg:           &sync.WaitGroup{},
    }
}
