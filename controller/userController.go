package controller

import (
	"context"
	"fmt"
	"go-jwt/database"
	"go-jwt/helper"
	"go-jwt/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(pass)
}

func VerifyPassword(userPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user model.User

		if err := c.BindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validateErr := validate.Struct(user)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})

		}
		password := HashPassword(*user.Password)
		user.Password = &password
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for phone number"})
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})

		}
		user.Created_at = time.Now()
		user.Updated_at = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		insertionNumber, insertionErr := userCollection.InsertOne(ctx, user)
		if insertionErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		c.JSON(http.StatusOK, insertionNumber)
	}

}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user model.User
		var foundUser model.User
		if err := c.BindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		token,refreshToken,_:=helper.GenerateAllTokens(*foundUser.Email,*foundUser.First_name,*foundUser.Last_name,*foundUser.User_type,foundUser.User_id)
		helper.UpdateAllTokens(token,refreshToken,foundUser.user_id)
		err = userCollection.FindOne(ctx,bson.M{"user_id":foundUser.User_id}).Decode(foundUser)
		if err!=nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusOK,foundUser)
	}
}

func GetUsers() {

}

// only adimin can access this function

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user model.User
		defer cancel()
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
