package config

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"gopkg.in/ini.v1"
)

const mySecret = "tualo-zw53htx6sX"

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		log.Println("Error: ", err)
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Decrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		log.Println("Error: ", err)
		return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

type ConfigurationClass struct {
	loaded   bool
	appID    string
	fileName string
	cfg      *ini.File
}

func (me *ConfigurationClass) SetAppID(id string) {
	me.appID = id
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func (me *ConfigurationClass) Load() {
	dirname, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}

	if !checkFileExists(path.Join(dirname, me.appID)) {
		os.Mkdir(path.Join(dirname, me.appID), os.ModePerm)
	}

	me.fileName = path.Join(dirname, me.appID, ".my.ini")

	log.Println("read data from ", me.fileName)
	if checkFileExists(me.fileName) {
		me.cfg, err = ini.Load(me.fileName)
		if err != nil {
			fmt.Printf("Fail to read file: %v", err)
			os.Exit(1)
		}
	} else {
		me.cfg = ini.Empty()
		me.cfg.SaveTo(me.fileName)
	}
	me.loaded = true
}

func (me *ConfigurationClass) Get(section string, key string) string {
	if me.loaded {
		if section == "credentials" && key == "password" {
			str, _ := Decrypt(me.cfg.Section(section).Key(key).String(), mySecret)
			return str
		}
		return me.cfg.Section(section).Key(key).String()
	}
	return ""
}
func (me *ConfigurationClass) Set(section string, key string, value string) {
	if me.loaded {
		if section == "credentials" && key == "password" {
			str, _ := Encrypt(value, mySecret)
			// log.Println("------",str,mySecret,value)
			me.cfg.Section(section).Key(key).SetValue(str)
		} else {
			me.cfg.Section(section).Key(key).SetValue(value)
		}
	}
}

func (me *ConfigurationClass) GetInt(section string, key string, defaultVal int) int {
	if me.loaded {
		return me.cfg.Section(section).Key(key).MustInt(defaultVal)
	}
	return 0
}
func (me *ConfigurationClass) SetInt(section string, key string, value int) {
	if me.loaded {
		me.cfg.Section(section).Key(key).SetValue(fmt.Sprintf("%d", value))
	}
}

func (me *ConfigurationClass) GetFloat64(section string, key string, defaultVal float64) float64 {
	if me.loaded {
		return me.cfg.Section(section).Key(key).MustFloat64(defaultVal)
	}
	return 0
}
func (me *ConfigurationClass) SetFloat64(section string, key string, value float64) {
	if me.loaded {
		me.cfg.Section(section).Key(key).SetValue(fmt.Sprintf("%f", value))
	}
}

func (me *ConfigurationClass) GetFloat32(section string, key string, defaultVal float32) float32 {
	if me.loaded {
		return float32(me.cfg.Section(section).Key(key).MustFloat64(float64(defaultVal)))
	}
	return 0
}
func (me *ConfigurationClass) SetFloat32(section string, key string, value float32) {
	if me.loaded {
		me.cfg.Section(section).Key(key).SetValue(fmt.Sprintf("%f", value))
	}
}

func (me *ConfigurationClass) GetBool(section string, key string, defaultVal bool) bool {
	if me.loaded {
		return me.cfg.Section(section).Key(key).MustBool(defaultVal)
	}
	return false
}
func (me *ConfigurationClass) SetBool(section string, key string, value bool) {
	if me.loaded {
		me.cfg.Section(section).Key(key).SetValue(fmt.Sprintf("%t", value))
	}
}

func (me *ConfigurationClass) Save() {
	if me.loaded {
		me.cfg.SaveTo(me.fileName)
	}
}

var o = &ConfigurationClass{
	loaded: false,
}

func Configuration() *ConfigurationClass {
	return o
}
