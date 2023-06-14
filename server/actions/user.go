package actions

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/davidalvarez305/review_poster/server/database"
	"github.com/davidalvarez305/review_poster/server/models"
	"github.com/davidalvarez305/review_poster/server/sessions"
	"github.com/davidalvarez305/review_poster/server/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func CreateUser(user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	_, err = generateToken(user.ID)

	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return database.DB.Save(&user).First(&user).Error
}

func UpdateUser(body models.User) (models.User, error) {
	var user models.User

	err := database.DB.Where("username = ?", body.Username).First(&user).Error

	user.Username = body.Username
	user.Email = body.Email

	_, err = generateToken(user.ID)

	if err != nil {
		return user, err
	}

	err = database.DB.Save(&user).First(&user).Error

	return user, err
}

func GetUserById(userId string) (models.User, error) {
	var user models.User

	err := database.DB.Where("id = ?", userId).First(&user).Error

	return user, err
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

	return fmt.Sprintf("%v", uId), nil
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

	user, err = GetUserById(uId)

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

	csrfToken, err := generateToken(user.ID)

	if err != nil {
		return err
	}

	sess.Set("csrf_token", csrfToken.UUID)

	err = sess.Save()

	return err
}

func RequestChangePasswordCode(user models.User) error {
	var token models.Token

	token, err := generateToken(user.ID)

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

	_, err = generateToken(user.ID)

	if err != nil {
		return newUser, err
	}

	newUser = user

	err = database.DB.Save(&newUser).First(&newUser).Error

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
	token, err := refreshAuthToken()

	if err != nil {
		return err
	}

	// Initialize Client & Service With Credentials
	client := config.Client(context.Background(), &token)

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

func GetCsrfTokenFromSession(c *fiber.Ctx) (string, error) {
	var token string

	sess, err := sessions.Sessions.Get(c)

	if err != nil {
		return token, err
	}

	csrf_token := sess.Get("csrf_token")

	if csrf_token == nil {
		return token, errors.New("token not found")
	}

	token = fmt.Sprintf("%v", csrf_token)

	return token, nil
}
