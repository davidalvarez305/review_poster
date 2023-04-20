package actions

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/davidalvarez305/review_poster/cms/server/database"
	"github.com/davidalvarez305/review_poster/cms/server/models"
	"github.com/davidalvarez305/review_poster/cms/server/sessions"
	"github.com/davidalvarez305/review_poster/cms/server/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func SaveUser(user models.User) error {
	return database.DB.Save(&user).First(&user).Error
}

func Logout(c *fiber.Ctx) error {
	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return err
	}

	err = sess.Destroy()

	return err
}

func Delete(user models.User) error {
	return database.DB.Delete(&user).Error
}

func CreateUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	token, err := GenerateToken()

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.Token = &token

	return SaveUser(user)
}

func UpdateUser(user models.User, body models.User) error {

	user.Username = body.Username
	user.Email = body.Email

	token, err := GenerateToken()

	if err != nil {
		return err
	}

	user.Token = &token

	return SaveUser(user)
}

func GetUserById(user models.User, userId string) error {
	return database.DB.Where("id = ?", userId).First(&user).Error
}

func GetUserIdFromSession(c *fiber.Ctx) (string, error) {
	var userId string
	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return userId, err
	}

	uId := sess.Get("userId")

	if uId == nil {
		return userId, errors.New("user not found")
	}

	userId = fmt.Sprintf("%v", uId)

	return userId, nil
}

func GetUserFromSession(c *fiber.Ctx) (models.User, error) {
	var user models.User

	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return user, err
	}

	userId := sess.Get("userId")

	if userId == nil {
		return user, errors.New("user not found")
	}

	uId := fmt.Sprintf("%v", userId)

	err = GetUserById(user, uId)

	if err != nil {
		return user, err
	}

	return user, nil
}

func Login(user models.User, c *fiber.Ctx) error {
	userPassword := user.Password
	err := database.DB.Where("username = ?", user.Username).First(&user).Error

	if err != nil {
		return errors.New("incorrect username")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword))

	if err != nil {
		return errors.New("incorrect password")
	}

	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return err
	}

	sess.Set("userId", user.ID)

	err = sess.Save()

	return err
}

func RequestChangePasswordCode(user models.User) error {
	var token models.Token

	token, err := GenerateToken()

	if err != nil {
		return err
	}

	return SendGmail(user, token.UUID)
}

func ChangePassword(user models.User, password string) (models.User, error) {
	var newUser models.User

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return newUser, err
	}

	user.Password = string(hashedPassword)

	token, err := GenerateToken()

	if err != nil {
		return newUser, err
	}

	user.Token = &token

	newUser = user

	err = SaveUser(newUser)

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func SendGmail(user models.User, uuidCode string) error {

	// Load & Read Credentials From Credentials JSON File
	ctx := context.Background()
	googlePath := os.Getenv("GOOGLE_JSON_PATH")
	path, err := utils.ResolveServerPath()

	if err != nil {
		return err
	}

	b, err := os.ReadFile(path + "/" + googlePath)

	if err != nil {
		return err
	}

	// Create OAuth2 Pointer Config Struct
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return err
	}

	// Refresh Token
	token, err := utils.RefreshAuthToken()

	if err != nil {
		return err
	}

	// Initialize Client & Service With Credentials
	client := config.Client(context.Background(), token)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		return err
	}

	// Craft & Send Message
	from := os.Getenv("GMAIL_EMAIL")
	to := user.Email
	title := "Change Password Request"
	message := fmt.Sprintf("Click to change your password: %s", os.Getenv("CONTENT_CLIENT_URL")+"/token/"+uuidCode)

	msgStr := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, title, message)

	msg := []byte(msgStr)

	gMessage := &gmail.Message{Raw: base64.URLEncoding.EncodeToString(msg)}

	_, err = srv.Users.Messages.Send("me", gMessage).Do()

	if err != nil {
		return err
	}

	return nil
}
