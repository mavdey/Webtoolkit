package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

func sendPost(c *fiber.Ctx, url string, body map[string]interface{}) []byte {
	logrus.Debug("sendPost")
	jsonBody, err := json.Marshal(body)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}

	// Отправляем POST-запрос на целевой URL
	targetURL := viper.GetString("backend.host") + url // Замените на свой URL
	response, err := http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	defer response.Body.Close()

	// Читаем ответ
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	logrus.Debug("targetURL: ", targetURL)
	logrus.Debug("body: ", body)
	logrus.Debug("response: ", response)
	fmt.Println("Response =", string(responseBody))
	return responseBody
}

func sendGet(c *fiber.Ctx, url string) []byte {
	logrus.Debug("sendGet")
	targetURL := viper.GetString("backend.host") + url
	response, err := http.Get(targetURL)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	defer response.Body.Close()

	// Читаем ответ
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	logrus.Debug("targetURL: ", targetURL)
	logrus.Debug("response: ", response)
	fmt.Println("Response =", string(responseBody))
	return responseBody
}

func sendRequest(c *fiber.Ctx, args ...interface{}) []byte {
	logrus.Debug("sendRequest")
	var strArgs []string
	var response *http.Response
	var err error
	var body map[string]interface{}

	for _, arg := range args {
		switch v := arg.(type) {
		case map[string]interface{}:
			body = v
		case string:
			strArgs = append(strArgs, v)
		}
		fmt.Printf("Type: %T, Value: %v\n", arg, arg)
	}
	targetURL := viper.GetString("backend.host") + strArgs[1]
	if strArgs[0] == "Get" {
		response, err = http.Get(targetURL)
	} else if strArgs[0] == "Post" {
		logrus.Debug("body: ", body)
		jsonBody, err := json.Marshal(body)
		if err != nil {
			c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
			return nil
		}
		response, err = http.Post(targetURL, "application/json", bytes.NewBuffer(jsonBody))
	}
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.Status(http.StatusInternalServerError).SendString("Internal Server Error")
		return nil
	}
	logrus.Debug("targetURL: ", targetURL)
	logrus.Debug("response: ", response)
	fmt.Println("Response =", string(responseBody))
	return responseBody
}
