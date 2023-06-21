package goexecutor

import (
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

func (c *goroutineControl) Work(param map[string]any, f func(param map[string]any)) {
    c.Add()
    go func(c *goroutineControl, param map[string]any, fu func(param map[string]any)) {
        defer c.Done()
        fu(param)
    }(c, param, f)
}

func (c *goroutineControl) Add() {
    c.wg.Add(1)
    c.GoroutineNum <- struct{}{}
    if maxGoroutineNum > 0 {
        maxGoroutine <- struct{}{}
    }
}

func (c *goroutineControl) GetCurrentGoroutineCount() int {
    return len(c.GoroutineNum)
}

func (c *goroutineControl) Done() {
    <-c.GoroutineNum
    if maxGoroutineNum > 0 {
        <-maxGoroutine
    }
    c.wg.Done()
}

func (c *goroutineControl) Wait() {
    c.wg.Wait()
}

func New(goroutineNum uint16) *goroutineControl {
    return &goroutineControl{
        GoroutineNum: make(chan struct{}, goroutineNum),
        wg:           &sync.WaitGroup{},
    }
}
