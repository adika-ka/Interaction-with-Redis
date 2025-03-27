package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
)

type Cacher interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

type cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cacher {
	return &cache{
		client: client,
	}
}

type User struct {
	ID   int
	Name string
	Age  int
}

func (c *cache) Set(key string, value interface{}) error {
	valueStr, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error write json: %w", err)
	}

	err = c.client.Set(key, string(valueStr), 0).Err()
	if err != nil {
		return fmt.Errorf("error setting value in Redis: %w", err)
	}
	return nil
}

func (c *cache) Get(key string) (interface{}, error) {
	res, err := c.client.Get(key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("not found by key %s", key)
	} else if err != nil {
		return nil, fmt.Errorf("ошибка получения данных из кэша: %w", err)
	}
	return res, nil
}

func main() {
	// Создание клиента Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer client.Close()

	cache := NewCache(client)

	// Установка значения по ключу
	err := cache.Set("some:key", "value")
	if err != nil {
		panic(err)
	}

	// Получение значения по ключу
	value, err := cache.Get("some:key")
	if err != nil {
		panic(err)
	}

	fmt.Println(value)

	user := &User{
		ID:   1,
		Name: "John",
		Age:  30,
	}
	// Установка значения по ключу
	err = cache.Set(fmt.Sprintf("user:%v", user.ID), user)
	if err != nil {
		panic(err)
	}

	// Получение значения по ключу
	value, err = cache.Get("user:1")
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
}
