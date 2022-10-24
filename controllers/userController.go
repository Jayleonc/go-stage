package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Jayleonc/go-stage/database"
	helper "github.com/Jayleonc/go-stage/helpers"
	"github.com/Jayleonc/go-stage/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"time"
)

// validator 用于对数据进行校验
var validate = validator.New()

var db = database.DBInstance()

func HashPassword(password string) string {
	// 使用密码生成哈希
	ins := sha256.New()
	ins.Write([]byte(password))
	result := ins.Sum([]byte(""))
	return hex.EncodeToString(result)
}

func VerifyPassword(userPassword string, dataPassword string) (bool, string) {
	// 使用用户输入的密码生成sha256 hash
	password := HashPassword(userPassword)
	check := true
	msg := ""
	if password != dataPassword {
		check = false
		msg = fmt.Sprintf("password of email is incorrect")
	}
	return check, msg
}

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// 根据要注册的用户 email 和 phone 查询数据库是否已经存在该用户数据
		find := db.Where("(email, phone) IN ?", [][]interface{}{{user.Email, user.Phone}}).Find(&user)
		if find.RowsAffected > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone already exists"})
			return
		}

		// 构建用户结构
		user.ID = uuid.NewString()
		user.UserId = user.ID
		// 将密码计算为 hash 值写入数据库
		password := HashPassword(user.Password)
		user.Password = password
		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		// 计算token
		token, refreshToken, _ := helper.GenerateAllTokens(user.Email, user.UserName, user.UserType, user.UserId)
		user.Token = token
		user.RefreshToken = refreshToken

		create := db.Debug().Create(&user)
		if create.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":  "注册成功",
			"code": 200,
			"data": &user,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 绑定json BindJSON()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. 使用 email 查找用户数据，不存在即报错返回
		var found models.User
		find := db.Where("email = ?", user.Email).Table("user").Find(&found)
		if find.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查找用户数据失败"})
			return
		}

		// 3. VerifyPassword 验证密码是否正确
		passwordIsValid, msg := VerifyPassword(user.Password, found.Password)
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"message": msg})
			return
		}

		// 4. 生成用户 token，并更新到数据库
		tokens, refreshToken, _ := helper.GenerateAllTokens(found.Email, found.UserName, found.UserType, found.UserId)
		// todo 如何确保用户数据完整性，判断用户数据是否被非法改动
		helper.UpdateAllTokens(tokens, refreshToken, found.UserId)

		//5. 登陆成功，返回用户数据
		find = db.Where("user_id = ?", found.UserId).Table("user").Find(&found)
		if find.Error != nil && find.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登陆失败"})
			return
		}
		c.JSON(http.StatusOK, found)
	}
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1 这个 api 只有管理员 ADMIN 才能访问，USER 访问，报错返回
		var userDto models.UserPageDto
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if c.ShouldBindQuery(&userDto) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "参数错误"})
		}
		// 构造查询 query 参数
		pageNum := userDto.PageNum
		pageSize := userDto.PageSize
		userName := userDto.UserName

		if pageNum <= 0 {
			pageNum = 1
		}

		var allUser []models.User
		// 分页查询
		find := db.Debug().Where("user_name LIKE ? ", "%"+userName+"%").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&allUser)
		if find.RowsAffected == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		}
		// 查询总数
		var total int64
		db.Table("user").Count(&total)
		pageNum = int(total / (int64(pageSize)))
		if total%int64(pageSize) != 0 {
			pageNum++
		}

		c.JSON(http.StatusOK, gin.H{
			"code":  200,
			"total": total,
			"data":  allUser,
		})

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
