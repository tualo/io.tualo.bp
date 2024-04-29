package config
import (
	"os"
	"log"
	"fmt"
	"errors"
	"gopkg.in/ini.v1"
	"path"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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
		log.Println("Error: ",err)
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
		log.Println("Error: ",err)
	 return "", err
	}
	cipherText := Decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

type ConfigurationClass struct {
	loaded bool
	appID string
	fileName string
	cfg *ini.File
}
func (this *ConfigurationClass) SetAppID(id string) {
	this.appID = id
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func (this *ConfigurationClass) Load() {
	dirname, err := os.UserConfigDir()
    if err != nil {
        log.Fatal( err )
    }

	if checkFileExists(path.Join(dirname,this.appID)) == false {
		os.Mkdir(path.Join(dirname,this.appID), os.ModePerm)
	}

	this.fileName = path.Join(dirname,this.appID,".my.ini")

	
	
	if checkFileExists(this.fileName) {
		this.cfg, err = ini.Load(this.fileName)
		if err != nil {
			fmt.Printf("Fail to read file: %v", err)
			os.Exit(1)
		}
	}else{
		this.cfg = ini.Empty()
		this.cfg.SaveTo(this.fileName)
	}
	this.loaded = true
}
func (this *ConfigurationClass) Get(section string, key string ) string {
	if this.loaded {
		if section == "credentials" && key == "password" {
			str,_:=Decrypt(this.cfg.Section(section).Key(key).String(),mySecret)
			return str
		}
		return this.cfg.Section(section).Key(key).String()
	}
	return ""
}
func (this *ConfigurationClass) Set(section string, key string,value string) {
	if this.loaded {
		if section == "credentials" && key == "password" {
			str,_:=Encrypt(value,mySecret)
			log.Println("------",str,mySecret,value)
			this.cfg.Section(section).Key(key).SetValue(str)
		}else{
			this.cfg.Section(section).Key(key).SetValue(value)
		}
	}
}

func (this *ConfigurationClass) GetInt(section string, key string, defaultVal int) int {
	if this.loaded {
		return this.cfg.Section(section).Key(key).MustInt(defaultVal)
	}
	return 0
}
func (this *ConfigurationClass) SetInt(section string, key string,value int) {
	if this.loaded {
		this.cfg.Section(section).Key(key).SetValue(fmt.Sprintf("%d",value))
	}
}

func (this *ConfigurationClass) GetFloat64(section string, key string,defaultVal float64) float64 {
	if this.loaded {
		return this.cfg.Section(section).Key(key).MustFloat64(defaultVal)
	}
	return 0
}
func (this *ConfigurationClass) SetFloat64(section string, key string,value float64) {
	if this.loaded {
		this.cfg.Section(section).Key(key).SetValue(fmt.Sprintf("%f",value))
	}
}

func (this *ConfigurationClass) GetFloat32(section string, key string, defaultVal float32) float32 {
	if this.loaded {
		return float32( this.cfg.Section(section).Key(key).MustFloat64(float64(defaultVal)) )
	}
	return 0
}
func (this *ConfigurationClass) SetFloat32(section string, key string,value float32) {
	if this.loaded {
		this.cfg.Section(section).Key(key).SetValue(fmt.Sprintf("%f",value))
	}
}

func (this *ConfigurationClass) GetBool(section string, key string,defaultVal bool) bool {
	if this.loaded {
		return this.cfg.Section(section).Key(key).MustBool(defaultVal)
	}
	return false
}
func (this *ConfigurationClass) SetBool(section string, key string,value bool) {
	if this.loaded {
		this.cfg.Section(section).Key(key).SetValue(fmt.Sprintf("%t",value))
	}
}

func (this *ConfigurationClass) Save() {
	if this.loaded {
		this.cfg.SaveTo(this.fileName)
	}
}
func NewConfigurationClass() *ConfigurationClass {
	o := &ConfigurationClass{
		loaded: false,
	}
	// o.SetPlayState( false )
	return o
}