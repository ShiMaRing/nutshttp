package nutshttp

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Product struct {
	Keyword  string
	Name     string
	Price    float64
	Date     time.Time
	ShopName string
	ImgURL   string
	Comment  string
}

const XANADU_BUCKET = "xanadu"

func (s *NutsHTTPServer) Search(context *gin.Context) {
	//获取用户搜索关键词
	keyword := context.Param("keyword")
	//尝试先从数据库获取
	cache, err := s.core.Get(XANADU_BUCKET, keyword)
	if err != nil {
		//说明爬取不到
		products := s.crawSkuFromWeb(keyword)
		if len(products) == 0 {
			WriteError(context, APIMessage{
				Message: "未找到相关商品",
			})
			return
		}
		//存入数据库
		productsJson, err := json.Marshal(products)
		if err != nil {
			log.Fatal(err)
		}
		err = s.core.Update(XANADU_BUCKET, keyword, string(productsJson), -1)
		if err != nil {
			log.Fatal(err)
		}
		WriteSucc(context, products)
	}
	var products = make([]Product, 0)
	err = json.Unmarshal([]byte(cache), &products)
	if err != nil {
		log.Fatal(err)
	}
	WriteSucc(context, products)
}

func (s *NutsHTTPServer) crawSkuFromWeb(keyword string) []Product {
	url := "https://search.dangdang.com/?key=" + keyword + "&show=list&act=input"
	c := colly.NewCollector()
	var products []Product
	c.OnHTML("#search_nature_rg ul li", func(e *colly.HTMLElement) {
		skuName := e.ChildText(".name a")
		skuPriceStr := e.ChildText(".price .search_now_price")
		skuPrice, _ := strconv.ParseFloat(strings.TrimPrefix(skuPriceStr, "¥"), 64)
		image := e.ChildAttr(".pic img", "data-original")
		if image == "" {
			return
		}
		skuImgURL := "https:" + image
		shopName := e.ChildText(".link a")
		if shopName == "" {
			shopName = "自营店铺"
		}
		comment := e.ChildText(".search_hot_word")

		//去掉skuName如 【*】 正则表达式标签
		index := strings.IndexRune(skuName, '】')
		if index != -1 {
			skuName = skuName[index+3:]
		}
		product := Product{
			Keyword:  keyword,
			Name:     skuName,
			Price:    skuPrice,
			Date:     time.Now(),
			ShopName: shopName,
			ImgURL:   skuImgURL,
			Comment:  comment,
		}
		products = append(products, product)
	})
	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
	c.Wait()
	return products
}
