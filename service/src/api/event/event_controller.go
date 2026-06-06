package event

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"512b.it/daytrack/src/api/middleware"
	"512b.it/daytrack/src/database"
	"512b.it/daytrack/src/models"
	"512b.it/daytrack/src/utils"
	"github.com/gin-gonic/gin"
)

type EventController struct {
	unauthenticatedRoute *gin.RouterGroup
	authenticatedRoute   *gin.RouterGroup

	db            *database.Database
	configuration models.Configuration
	event         *Service
	logger        utils.ContextLogger
}

func Inject(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	db *database.Database,
	event *Service,
	configuration models.Configuration,
) {
	controller := new(unauthenticatedRoute, authenticatedRoute, db, event, configuration)
	controller.injectUnauthenticatedRoutes()
	controller.injectAuthenticatedRoutes()
}

func new(
	unauthenticatedRoute *gin.RouterGroup,
	authenticatedRoute *gin.RouterGroup,
	db *database.Database,
	event *Service,
	configuration models.Configuration,
) *EventController {
	logger := utils.InitServiceLogger("EventController")

	return &EventController{
		unauthenticatedRoute: unauthenticatedRoute,
		authenticatedRoute:   authenticatedRoute,
		db:                   db,
		event:                event,
		configuration:        configuration,
		logger:               logger,
	}
}

func (c *EventController) injectUnauthenticatedRoutes() {
	v1 := c.unauthenticatedRoute.Group("v1")
	v1.Use(middleware.AuthApiKeyGuards(c.configuration, c.db))
	{
		v1.POST("/events/:username/:track_name", utils.RateLimit(c.configuration.RateLimitBurst, c.configuration.RateLimitInterval), c.createEventRoute())
		v1.GET("/events/:username/:track_name", c.listEventRoute())
	}
}

func (c *EventController) injectAuthenticatedRoutes() {}

// @Tags event
// @Security ApiKeyAuth
// @Schemes https
// @Router /v1/events/{username}/{track_name} [POST]
// @Summary Create an event
// @Description Create an event for the authenticated user
// @Param username path string true "Username"
// @Param track_name path string true "Track name"
// @Param key query string false "API Key"
// @Param created_at query string false "Created at"
// @Param quantity query int false "Quantity"
// @Accept json
// @Produce json
// @Success 201 {object} models.ResponseModel
func (c *EventController) createEventRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var createdAt *time.Time

		username := ctx.Param("username")
		trackName := ctx.Param("track_name")
		createdAtString := ctx.DefaultQuery("created_at", "")
		quantity, err := strconv.Atoi(ctx.DefaultQuery("quantity", "1"))

		if err != nil {
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid quantity"))
			return
		}

		if createdAtString != "" {
			createdAtx, err := time.Parse(time.RFC3339, createdAtString)
			if err != nil {
				c.logger(ctx).Err(err).Msg("Failed to parse created_at")
				ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid after parameter, expected RFC 3339 format"))
				return
			}
			createdAt = &createdAtx
		}

		var uid *int64
		if rawID, authErr := utils.GetAuthenticatedUserID(ctx); authErr == nil {
			uid = &rawID
		}

		if err = c.event.CreateEvent(uid, username, trackName, quantity, createdAt); err != nil {
			c.logger(ctx).Err(err).Msg("Failed to create event")
			if errors.Is(err, ErrorEventCannotBeCalledByYou) || errors.Is(err, database.ErrorNotFound) {
				ctx.JSON(403, models.NewError(models.ErrorForbidden, "event cannot be called by you"))
				return
			}
			ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to create event"))
			return
		}

		ctx.JSON(201, models.NewSuccess("event registered", ""))
	}
}

// @Tags event
// @Security ApiKeyAuth
// @Schemes https
// @Router /v1/events/{username}/{track_name} [GET]
// @Summary List events
// @Description List all events for a track
// @Param username path string true "Username"
// @Param track_name path string true "Track name"
// @Param key query string false "API Key"
// @Param after query string false "List events after RFC 3339 timestamp"
// @Param list_by query string false "Aggregation (day, month, raw)" Enums(day, month, raw)
// @Param format query string false "Response format (json or html)" Enums(json, html)
// @Accept json
// @Produce json,text/html
// @Success 200 {object} models.List[models.Day]
func (c *EventController) listEventRoute() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var events []models.Day
		var after *time.Time

		username := ctx.Param("username")
		trackName := ctx.Param("track_name")
		afterString := utils.EmptyIsNull(ctx.DefaultQuery("after", ""))

		listBy := models.ListBy(ctx.DefaultQuery("list_by", string(models.ListByDay)))
		format := models.Format(ctx.DefaultQuery("format", string(models.FormatJSON)))

		var uid *int64
		if rawID, authErr := utils.GetAuthenticatedUserID(ctx); authErr == nil {
			uid = &rawID
		}

		if afterString != nil {
			afterX, err := time.Parse(time.RFC3339, *afterString)
			if err != nil {
				ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid after parameter, expected RFC 3339 format"))
				return
			}
			after = &afterX
		}

		queryListBy := listBy
		if format == models.FormatHTML {
			queryListBy = models.ListByDay
		}

		switch queryListBy {
		case models.ListByDay:
			events, err = c.event.ListEventsByDays(uid, username, trackName, after)
		case models.ListByMonth:
			events, err = c.event.ListEventsByMonth(uid, username, trackName, after)
		case models.ListByRaw:
			events, err = c.event.ListEvents(uid, username, trackName, after)
		default:
			ctx.JSON(400, models.NewError(models.ErrorBadRequest, "invalid list_by"))
			return
		}

		if err != nil {
			c.logger(ctx).Err(err).Msg("Failed to list events")
			if errors.Is(err, ErrorEventCannotBeCalledByYou) || errors.Is(err, database.ErrorNotFound) {
				ctx.JSON(403, models.NewError(models.ErrorForbidden, "event cannot be called by you"))
			} else {
				ctx.JSON(500, models.NewError(models.ErrorInternalServerError, "failed to list events"))
			}
			return
		}

		switch format {
		case models.FormatHTML:
			c.renderHeatmap(ctx, username, trackName, events)
		default:
			ctx.JSON(200, models.List[models.Day]{
				Items: events,
			})
		}
	}
}

func (c *EventController) renderHeatmap(ctx *gin.Context, username, trackName string, events []models.Day) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("X-Content-Type-Options", "nosniff")

	byDate := make(map[string]int)
	for _, e := range events {
		key := e.Date.Format("2006-01-02")
		byDate[key] += e.Quantity
	}

	maxQty := 1
	for _, q := range byDate {
		if q > maxQty {
			maxQty = q
		}
	}

	today := time.Now()
	weekday := int(today.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	end := today.AddDate(0, 0, 7-weekday)
	start := end.AddDate(0, 0, -20*7+1)

	type cell struct {
		date string
		qty  int
	}
	var allCells []cell
	cursor := start
	for !cursor.After(end) {
		key := cursor.Format("2006-01-02")
		q := byDate[key]
		allCells = append(allCells, cell{date: key, qty: q})
		cursor = cursor.AddDate(0, 0, 1)
	}

	const cellsPerWeek = 7
	type week struct {
		cells []cell
	}
	var weeks []week
	for i := 0; i < len(allCells); i += cellsPerWeek {
		end := i + cellsPerWeek
		if end > len(allCells) {
			end = len(allCells)
		}
		weeks = append(weeks, week{cells: allCells[i:end]})
	}

	monthNames := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	type monthLabel struct {
		index int
		label string
	}
	var monthLabels []monthLabel
	lastMonth := -1
	for wi, w := range weeks {
		if len(w.cells) == 0 {
			continue
		}
		parsed, _ := time.Parse("2006-01-02", w.cells[0].date)
		m := int(parsed.Month())
		if m != lastMonth {
			monthLabels = append(monthLabels, monthLabel{index: wi, label: monthNames[m-1]})
			lastMonth = m
		}
	}

	dayLabels := []string{"", "Mon", "", "Wed", "", "Fri", ""}

	var cellsHTML strings.Builder
	for _, w := range weeks {
		cellsHTML.WriteString(`<div class="hw">`)
		for _, c := range w.cells {
			level := 0
			if c.qty > 0 {
				level = (c.qty * 4) / maxQty
				if level > 4 {
					level = 4
				}
			}
			cls := "hc"
			if level == 0 {
				cls += " e"
			}
			cls += fmt.Sprintf(" l-%d", level)
			cellsHTML.WriteString(fmt.Sprintf(`<div class="%s" title="%s: %d"></div>`, cls, c.date, c.qty))
		}
		cellsHTML.WriteString("</div>")
	}

	var monthLabelsHTML strings.Builder
	for _, ml := range monthLabels {
		left := ml.index * 16
		monthLabelsHTML.WriteString(fmt.Sprintf(`<span class="hml" style="left:%dpx">%s</span>`, left, ml.label))
	}

	var dayLabelsHTML strings.Builder
	for _, dl := range dayLabels {
		if dl == "" {
			dayLabelsHTML.WriteString(`<div class="hdl"></div>`)
		} else {
			dayLabelsHTML.WriteString(fmt.Sprintf(`<div class="hdl">%s</div>`, dl))
		}
	}

	eu := htmlEsc(username)
	et := htmlEsc(trackName)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s / %s — Activity</title>
<meta name="description" content="Activity heatmap for %s/%s">
<meta name="robots" content="noindex">
<link rel="icon" href="/logo.svg">
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Oxygen,Ubuntu,Cantarell,sans-serif;background:#0f172a;color:#e2e8f0;padding:24px}
.hg{display:flex;gap:6px}
.hl{display:grid;grid-template-rows:14px repeat(7,13px);gap:3px;margin-right:2px}
.hsp{height:14px}
.hdl{font-size:10px;color:#64748b;line-height:13px;height:13px}
.hws{overflow-x:auto}
.hmr{position:relative;height:14px;margin-bottom:2px;min-width:0}
.hml{position:absolute;top:0;font-size:10px;color:#64748b;white-space:nowrap}
.hcs{display:flex;gap:3px}
.hw{display:grid;grid-template-rows:repeat(7,13px);gap:3px}
.hc{width:13px;height:13px;border-radius:3px;cursor:default}
.hc.e,.hc.l-0{outline:1px solid #475569;outline-offset:-1px;background:transparent}
.hc.l-1{background:#a78bfa}
.hc.l-2{background:#8b5cf6}
.hc.l-3{background:#6d28d9}
.hc.l-4{background:#4c1d95}
.lg{display:flex;align-items:center;gap:3px;margin-top:8px;justify-content:flex-end;font-size:10px;color:#64748b}
.lg .hc{width:0;min-width:11px;height:11px}
</style>
</head>
<body>
<div class="hg">
<div class="hl">
<div class="hsp"></div>
%s
</div>
<div class="hws">
<div class="hmr">%s</div>
<div class="hcs">%s</div>
</div>
</div>
<div class="lg">
<span>Less</span>
<div class="hc l-0"></div>
<div class="hc l-1"></div>
<div class="hc l-2"></div>
<div class="hc l-3"></div>
<div class="hc l-4"></div>
<span>More</span>
</div>
</body>
</html>`,
		et, eu, eu, et,
		dayLabelsHTML.String(),
		monthLabelsHTML.String(),
		cellsHTML.String())

	ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func htmlEsc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;", "'", "&#39;").Replace(s)
}
