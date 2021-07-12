// 基金

package routes

import (
	"net/http"
	"strings"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/x-stock/core"
	"github.com/axiaoxin-com/x-stock/datacenter"
	"github.com/axiaoxin-com/x-stock/models"
	"github.com/axiaoxin-com/x-stock/services"
	"github.com/axiaoxin-com/x-stock/version"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// ParamFundIndex FundIndex 请求参数
type ParamFundIndex struct {
	PageNum  int    `json:"page_num"  form:"page_num"`
	PageSize int    `json:"page_size" form:"page_size"`
	Sort     int    `json:"sort"      form:"sort"`
	Type     string `json:"type"      form:"type"`
}

// FundIndex godoc
func FundIndex(c *gin.Context) {
	fundList := services.Fund4433List
	p := ParamFundIndex{
		PageNum:  1,
		PageSize: 10,
		Sort:     models.FundSortTypeWeek,
	}
	if err := c.ShouldBind(&p); err != nil {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金",
			"Error":     err.Error(),
		}
		c.HTML(http.StatusOK, "fund_index.html", data)
		return
	}

	// 过滤
	if p.Type != "" {
		fundList = fundList.FilterByType(p.Type)
	}
	// 排序
	if p.Sort > 0 {
		fundList.Sort(models.FundSortType(p.Sort))
	}
	// 分页
	totalCount := len(fundList)
	pagi := goutils.PaginateByPageNumSize(totalCount, p.PageNum, p.PageSize)
	result := fundList[pagi.StartIndex:pagi.EndIndex]
	data := gin.H{
		"Env":           viper.GetString("env"),
		"Version":       version.Version,
		"PageTitle":     "X-STOCK | 基金",
		"URLPath":       "/fund",
		"FundList":      result,
		"Pagination":    pagi,
		"IndexParam":    p,
		"UpdatedAt":     services.SyncFundTime.Format("2006-01-02 15:04:05"),
		"AllFundCount":  len(services.FundAllList),
		"Fund4433Count": totalCount,
		"FundTypes":     services.Fund4433TypeList,
	}
	c.HTML(http.StatusOK, "fund_index.html", data)
	return
}

// ParamFundFilter FundFilter 请求参数
type ParamFundFilter struct {
	ParamFundListFilter models.ParamFundListFilter
	ParamFundIndex      ParamFundIndex
}

// FundFilter godoc
func FundFilter(c *gin.Context) {
	p := ParamFundFilter{
		ParamFundListFilter: models.ParamFundListFilter{
			MinScale:             2.0,
			MaxScale:             50.0,
			MinEstabYears:        5.0,
			MinManagerYears:      5.0,
			Year1RankRatio:       25.0,
			ThisYear235RankRatio: 25.0,
			Month6RankRatio:      33.33,
			Month3RankRatio:      33.33,
		},
		ParamFundIndex: ParamFundIndex{
			PageNum:  1,
			PageSize: 10,
			Sort:     0,
		},
	}
	if err := c.ShouldBind(&p); err != nil {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金 | 基金严选",
			"Error":     err.Error(),
		}
		c.HTML(http.StatusOK, "fund_filter.html", data)
		return
	}
	fundList := services.FundAllList.Filter(c, p.ParamFundListFilter)
	fundTypes := fundList.Types()
	// 过滤
	if p.ParamFundIndex.Type != "" {
		fundList = fundList.FilterByType(p.ParamFundIndex.Type)
	}
	// 排序
	if p.ParamFundIndex.Sort > 0 {
		fundList.Sort(models.FundSortType(p.ParamFundIndex.Sort))
	}
	// 分页
	pagi := goutils.PaginateByPageNumSize(len(fundList), p.ParamFundIndex.PageNum, p.ParamFundIndex.PageSize)
	result := fundList[pagi.StartIndex:pagi.EndIndex]
	data := gin.H{
		"Env":         viper.GetString("env"),
		"Version":     version.Version,
		"PageTitle":   "X-STOCK | 基金 | 基金严选",
		"URLPath":     "/fund/filter",
		"FundList":    result,
		"Pagination":  pagi,
		"IndexParam":  p.ParamFundIndex,
		"FilterParam": p.ParamFundListFilter,
		"FundTypes":   fundTypes,
	}
	c.HTML(http.StatusOK, "fund_filter.html", data)
	return
}

// ParamFundCheck FundCheck 请求参数
type ParamFundCheck struct {
	// 基金代码
	Code string `json:"fundcode"                 form:"fundcode"`
	// 基金规模最小值（亿）
	MinScale float64 `json:"min_scale"                form:"min_scale"`
	// 基金规模最大值（亿）
	MaxScale float64 `json:"max_scale"                form:"max_scale"`
	// 基金经理管理该基金最低年限
	MinManagerYears float64 `json:"min_manager_years"        form:"min_manager_years"`
	// 最近一年收益率排名比
	Year1RankRatio float64 `json:"year_1_rank_ratio"        form:"year_1_rank_ratio"`
	// 今年来、最近两年、最近三年、最近五年收益率排名比
	ThisYear235RankRatio float64 `json:"this_year_235_rank_ratio" form:"this_year_235_rank_ratio"`
	// 最近六月收益率排名比
	Month6RankRatio float64 `json:"month_6_rank_ratio"       form:"month_6_rank_ratio"`
	// 最近三月收益率排名比
	Month3RankRatio float64 `json:"month_3_rank_ratio"       form:"month_3_rank_ratio"`
	// 1,3,5年波动率平均值的最大值
	Max135AvgStddev float64 `json:"max_135_avg_stddev"       form:"max_135_avg_stddev"`
	// 1,3,5年夏普比率平均值的最小值
	Min135AvgSharp float64 `json:"min_135_avg_sharp"        form:"min_135_avg_sharp"`
	// 1,3,5年最大回撤率平均值的最大值
	Max135AvgRetr float64 `json:"max_135_avg_retr"         form:"max_135_avg_retr"`
	// 是否检测持仓个股
	CheckStocks bool `json:"check_stocks"             form:"check_stocks"`
	// 股票检测参数
	StockCheckerOptions core.CheckerOptions
}

// FundCheck godoc
func FundCheck(c *gin.Context) {
	p := ParamFundCheck{
		MinScale:             2.0,
		MaxScale:             50.0,
		MinManagerYears:      5.0,
		Year1RankRatio:       25.0,
		ThisYear235RankRatio: 25.0,
		Month6RankRatio:      33.33,
		Month3RankRatio:      33.33,
		Max135AvgStddev:      25.0,
		Min135AvgSharp:       1.0,
		Max135AvgRetr:        25.0,
		CheckStocks:          false,
	}
	if err := c.ShouldBind(&p); err != nil {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金 | 基金检测",
			"Error":     err.Error(),
		}
		c.JSON(http.StatusOK, data)
		return
	}

	if p.Code == "" {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金 | 基金检测",
			"Error":     "请填写基金代码",
		}
		c.JSON(http.StatusOK, data)
		return
	}

	fundresp, err := datacenter.EastMoney.QueryFundInfo(c, p.Code)
	if err != nil {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金 | 基金检测",
			"Error":     err.Error(),
		}
		c.JSON(http.StatusOK, data)
		return
	}
	fund := models.NewFund(c, fundresp)

	if !p.CheckStocks {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金 | 基金检测",
			"Fund":      fund,
			"Param":     p,
		}
		c.JSON(http.StatusOK, data)
		return
	}
	checker := core.NewChecker(c, p.StockCheckerOptions)
	checkResult, err := checker.CheckFundStocks(c, fund)
	if err != nil {
		data := gin.H{
			"Env":       viper.GetString("env"),
			"Version":   version.Version,
			"PageTitle": "X-STOCK | 基金 | 基金检测",
			"Error":     err.Error(),
		}
		c.JSON(http.StatusOK, data)
		return
	}
	data := gin.H{
		"Env":              viper.GetString("env"),
		"Version":          version.Version,
		"PageTitle":        "X-STOCK | 基金 | 基金检测",
		"Fund":             fund,
		"StockCheckResult": checkResult,
		"Param":            p,
	}
	c.JSON(http.StatusOK, data)
	return
}

// FundSimilarity 基金相似度
// http://localhost:1688/fund/similarity?codes=001975,163807,001869,519133,519644,270028,377530,550009,210003,002160,166301,001365,519642,000073,001808,001279,001397,000592,001938,008314,001679,163406,162605
// http://localhost:1688/fund/similarity?codes=001975,519133,519644,001365
// http://localhost:1688/fund/similarity?codes=001975,519133,519644,001365,550009
// http://localhost:1688/fund/similarity?codes=001975,519133,519644,001365,001279
// http://localhost:1688/fund/similarity?codes=001975,519133,519644,001365,166301
// http://localhost:1688/fund/similarity?codes=001975,519133,519644,001365,163406
// http://localhost:1688/fund/similarity?codes=001975,519133,519644,001365,163406,166301
// http://localhost:1688/fund/similarity?codes=001975,519133,519644,001365,163406,166301,001279,550009
func FundSimilarity(c *gin.Context) {
	codes := c.Query("codes")
	if codes == "" {
		return
	}
	codeList := strings.Split(codes, ",")
	checker := core.NewChecker(c, core.DefaultCheckerOptions)
	data, err := checker.GetFundStocksSimilarity(c, codeList)
	if err != nil {
		c.IndentedJSON(200, gin.H{"err": err.Error()})
		return
	}
	c.IndentedJSON(200, data)
	return
}
