package user

import (
	"strconv"
	"time"

	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	unauthenticatedRoute *gin.RouterGroup
	authenticatedRoute   *gin.RouterGroup

	configuration models.Configuration
	user          *Service
	logger        utils.ContextLogger
}

func Inject(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	user *Service,
	configuration models.Configuration,
) {
	controller := new(unauthenticatedRoute, authenticatedRoute, user, configuration)
	controller.injectUnauthenticatedRoutes()
	controller.injectAuthenticatedRoutes()
}

func new(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	user *Service,
	configuration models.Configuration,
) *UserController {
	logger := utils.InitServiceLogger("UserController")

	return &UserController{
		unauthenticatedRoute: unauthenticatedRoute,
		authenticatedRoute:   authenticatedRoute,
		user:                 user,
		configuration:        configuration,
		logger:               logger,
	}
}

func (c *UserController) injectUnauthenticatedRoutes() {
	v1 := c.unauthenticatedRoute.Group("v1")
	{
		v1.GET("users/challenge", c.createUserChallenge())
		v1.POST("users/signup", c.createUserRoute())
		v1.POST("users/signin", c.signinUserRoute())
	}
}

func (c *UserController) injectAuthenticatedRoutes() {
	// v1 := c.authenticatedRoute.Group("v1")
	// {
	// }
}

// @Tags user
// @Schemes https
// @Router /v1/users/challenge [GET]
// @Summary Create a user challenge
// @Description Create a user challenge given an email address
// @Accept json
// @Param email query string true "The email address of the user"
// @Produce json
// @Success 200 {object} models.Challenge
func (c *UserController) createUserChallenge() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.Query("email")

		challenge := strconv.FormatInt(utils.SecureRandom(int64(c.configuration.Opaque.ChallengeRange)), 10) + "-" + strconv.FormatInt(time.Now().Unix(), 10)

		salt := c.user.GenerateSalt(ctx, c.configuration.Opaque.SaltNonce, email)

		ctx.JSON(200, models.Challenge{
			Salt:      salt,
			Challenge: challenge,
		})
	}
}

// @Tags user
// @Schemes https
// @Router /v1/users/signup [POST]
// @Summary Create a user
// @Description Create a new user
// @Accept json
// @Param request body models.User true "the user body to create"
// @Produce json
// @Success 201 {object} models.ResponseModel
func (c *UserController) createUserRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var user models.User

		if err := ctx.ShouldBindJSON(&user); err != nil {
			c.logger(ctx).Err(err).Msg("Invalid request")
			ctx.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if _, err := c.user.CreateUser(ctx, user.Username, user.Email, user.PublicKey); err != nil {
			c.logger(ctx).Err(err).Msg("Unable to insert the user")
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		ctx.JSON(201, models.NewSuccess("User created", ""))
	}
}

// @Tags user
// @Schemes https
// @Router /v1/users/signin [POST]
// @Summary Sign in a user
// @Description Sign in a user
// @Accept json
// @Param request body  models.SignIn true "the signed challenge to enter the system"
// @Produce json
// @Success 201 {object} models.ResponseModel
func (c *UserController) signinUserRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var auth models.SignIn
		var token *models.Token

		if err := ctx.ShouldBindJSON(&auth); err != nil {
			c.logger(ctx).Err(err).Msg("Invalid request")
			ctx.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if token, err = c.user.SigninUser(ctx, auth.Email, auth.Challenge, auth.SignedChallenge); err != nil {
			c.logger(ctx).Err(err).Msg("Unable to sign in the user")
			ctx.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		ctx.JSON(200, token)
	}
}
