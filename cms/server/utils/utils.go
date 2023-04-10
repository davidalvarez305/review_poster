package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/davidalvarez305/review_poster/cms/server/models"
	"golang.org/x/crypto/bcrypt"
)

func GetIds(arr string) ([]int, error) {
	var ids []int
	idList := strings.Split(arr, ",")
	for _, i := range idList {
		j, err := strconv.Atoi(i)
		if err != nil {
			return ids, err
		}
		ids = append(ids, j)
	}
	return ids, nil
}

func GetWordId(words []models.Word, str string) int {
	var w int

	for i := 0; i < len(words); i++ {
		if words[i].Name == str {
			w = words[i].ID
			return w
		}
	}
	return -1
}

func GenerateAPIToken(str string) string {
	keys := str + fmt.Sprintf("%+v", time.Now().Unix())
	hash, err := bcrypt.GenerateFromPassword([]byte(keys), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(hash)
}

func ResolveServerPath() (string, error) {
	var p string

	if os.Getenv("PRODUCTION") == "1" {
		return ".", nil
	}

	u, err := user.Current()

	if err != nil {
		return p, err
	}

	p = u.HomeDir + os.Getenv("SERVER_PATH")

	return p, nil
}
