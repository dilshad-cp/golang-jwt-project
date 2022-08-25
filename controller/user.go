package controller

import (
	"context"
	"time"
	"github.com/dilshad-cp/golang-jwt-project/helpers"
	"github.com/dilshad-cp/golang-jwt-project/database"
	"github.com/dilshad-cp/golang-jwt-project/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate := validator.New()

func HashPassword(password string) string{
	bcrypt.GeneratePassword([]byte(password), 14)
	if err !={
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""
	if err != nil{
		msg := fmt.Sprintf("Email or password is incorrect")
		check = false
	}
	return check, msg
}

func Signup() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		var user models.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(user)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		defer cancel()
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password4
		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Phone})
		defer cancel()
		if err != nil{
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured"})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This phone or email already exists"})
		}
		user.Created_at, _ = time.Parse(time.RFC339, time.Now().Format(time.RFC339))
		user.Updated_at, _ = time.Parse(time.RFC339, time.Now().Format(time.RFC339))
		user.ID = primitive.NewObjectId()
		user.User_id = user.ID.Hex()

		token, refreshToken := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token := &token
		user.Refresh_token := &refreshToken

		resultInsertionNumber, insertErr = userCollection.InsertOne(ctx, user)
		if insertErr != nil{
			msg := fmt.Sprintf("User cannot be created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		}
		defer cancel()
		c.JSON(http.StatusOk, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if passwordIsValid != true{
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if foundUser.Email == nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, *foundUser.User_id)
		
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		
		c.JSON(http.StatusOk, foundUser)

	}
}
func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		helper.GetUserType(c, "ADMIN"); err != nil{
			json.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1{
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page <1{
			page = 1
		}

		startIndex := ( page - 1 ) * recordPerPage
		startIndex, err := strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match": bson.D{{}}}}
		groupStage := bson.D{{"$group": bson.D{
			{"_id": bson.D{{"_id": "null"}},
			{"total": bson.D{{"sum": 1}}},
			{"data": bson.D{{"$push": "$$ROOT"}}}
		}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id": 0},
					{"total_count": 1},
					{"user_items": bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				} 
			}
		}


		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage
		})

		defer cancel()

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error on listing items"})
		}
		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOk, allUsers[0])
	}

}
func GetUser() gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Param("user_id")
		

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)

		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOk, user)
	}
}
