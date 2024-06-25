package user

import (
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
)

//go:embed magic_link.html
var magicLinkTemplate string

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
		v1.POST("users/magic-link", c.createMagicLinkRoute())
		v1.GET("users/magic-link", c.getMagicLinkRoute())
	}
}

func (c *UserController) injectAuthenticatedRoutes() {
	v1 := c.authenticatedRoute.Group("v1")
	{
		v1.GET("users/self", c.getUserRoute())
	}
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
		var user models.CreateUser

		c.logger(ctx).Info().Msg("Creating a new user")

		if err := ctx.ShouldBindJSON(&user); err != nil {
			c.logger(ctx).Err(err).Msg("Invalid request")
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "Invalid request"))
			return
		}

		if !models.IsValidTrackName(user.Username) {
			c.logger(ctx).Error().Msgf("Invalid username %s", user.Username)
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "Invalid username"))
			return
		}

		if _, err := c.user.CreateUser(ctx, user.Username, user.Email, user.PublicKey); err != nil {
			c.logger(ctx).Err(err).Msg("Unable to insert the user")

			if errors.Is(err, database.ErrorDuplicate) {
				c.logger(ctx).Info().Msg("User not created, fail silently")
				ctx.JSON(201, models.NewSuccess("User created", ""))
				return
			}

			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "Unable to insert the user"))
			return
		}

		ctx.JSON(201, models.NewSuccess("User created", ""))

		c.logger(ctx).Info().Msg("User created")

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
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "Invalid request"))
			return
		}

		if token, err = c.user.SigninUser(ctx, auth.Email, auth.Challenge, auth.SignedChallenge); err != nil {
			c.logger(ctx).Err(err).Msg("Unable to sign in the user")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "Unable to sign in the user"))
			return
		}

		ctx.JSON(200, token)
	}
}

// @Tags user
// @Schemes https
// @Router /v1/users/magic-link [POST]
// @Summary Create a user
// @Description Create a new user
// @Accept json
// @Param request body models.User true "the user body to create"
// @Produce json
// @Success 201 {object} models.ResponseModel
func (c *UserController) createMagicLinkRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var user models.MagicLink

		c.logger(ctx).Info().Msg("Attempting to create a magic link")

		if err := ctx.ShouldBindJSON(&user); err != nil {
			c.logger(ctx).Err(err).Msg("Invalid request")
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "Invalid request"))
			return
		}

		if err = c.user.SendMagicLink(ctx, user.Email); err != nil {
			c.logger(ctx).Err(err).Msg("Unable to send the magic link")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "Unable to send the magic link"))
			return
		}

		ctx.JSON(201, models.NewSuccess("Magic link created", ""))

		c.logger(ctx).Info().Msg("Magic link created")
	}
}

// @Tags user
// @Schemes https
// @Router /v1/users/magic-link [GET]
// @Summary Create a user
// @Description Create a new user
// @Accept json
// @Param request body models.User true "the user body to create"
// @Produce json
// @Success 201 {object} models.ResponseModel
func (c *UserController) getMagicLinkRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.logger(ctx).Info().Msg("Attempting to create a magic link")

		key := ctx.Query("key")

		println(key)
		url := "dailytrack://"
		if c.configuration.Environment == "development" {
			url = "exp://127.0.0.1:8081/--/"
		}

		url = fmt.Sprintf("%sauthorize?key=%s", url, key)

		results := strings.ReplaceAll(magicLinkTemplate, "{{.Url}}", url)

		println(results)

		ctx.Data(200, "text/html; charset=utf-8", []byte(results))

		c.logger(ctx).Info().Msg("Magic link retrieved")
	}
}

// @Tags user
// @Schemes https
// @Router /v1/users/self [GET]
// @Summary Get user info
// @Description Get user info
// @Accept json
// @Produce json
// @Success 200 {object} models.User
func (c *UserController) getUserRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var err error
		var userID int64
		var user *models.User

		if userID, err = utils.GetAuthenticatedUserID(ctx); err != nil {
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to get authenticated user id"))
			return
		}

		c.logger(ctx).Info().Msgf("Attempting to get user info %d", userID)
		if user, err = c.user.GetUser(ctx, userID); err != nil {
			c.logger(ctx).Err(err).Msg("Unable to get the user")
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "Unable to get the user"))
			return
		}

		ctx.JSON(200, user)
	}
}
