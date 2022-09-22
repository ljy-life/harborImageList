package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// harbor 信息
type Infomation struct {
	Schema string
	URL    string
	User   string
	Pass   string
}

// 镜像 tag
type ImageTag struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// 获取 harbor 仓库的信息，返回列表
func (i *Infomation) HarborProject() ([]string, error) {
	var rr []map[string]interface{}
	u := i.Schema + "://" + i.URL + "/api/v2.0/projects"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		panic(err)
	}

	params := make(url.Values)
	params.Set("page_size", "100")
	req.URL.RawQuery = params.Encode()
	req.SetBasicAuth(i.User, i.Pass)
	r, err := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(r.Body)

	json.Unmarshal(body, &rr)

	arrayHarbor := make([]string, 0)
	for _, v := range rr {
		b, _ := json.Marshal(v)
		mapStr, err := simplejson.NewJson(b)
		if err != nil {
			panic(err)
		}
		pro_1, _ := mapStr.Get("name").String()
		arrayHarbor = append(arrayHarbor, pro_1)
		if err != nil {
			panic(err)
		}
	}
	return arrayHarbor, err
}

// 根据镜像名称获取镜像 project_id
func (i *Infomation) HarborID(v string) (int, error) {
	var id int
	u := i.Schema + "://" + i.URL + "/api/v2.0/projects/" + v
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		panic(err)
	}
	params := make(url.Values)
	req.URL.RawQuery = params.Encode()
	req.SetBasicAuth(i.User, i.Pass)
	r, err := http.DefaultClient.Do(req)
	body, err := ioutil.ReadAll(r.Body)

	var rr map[string]interface{}
	json.Unmarshal(body, &rr)

	b, _ := json.Marshal(rr)
	mapStr, err := simplejson.NewJson(b)
	if err != nil {
		panic(err)
	}
	id, err = mapStr.Get("project_id").Int()
	if err != nil {
		panic(err)
	}
	return id, err
}


// 根据提供的 harbor 仓库的信息，获取仓库中的镜像信息列表
func (i *Infomation) HarborImage(project []string) ([]string, error) {
	arrayimage := make([]string, 0)
	var err error
	for _, p := range project {
		PorjectID, _ := i.HarborID(p)
		pages := PorjectID / 100
		pages = pages + 1
		for page := 1; page <= pages; page++ {
			s := strconv.Itoa(page)
			u := i.Schema + "://" + i.URL + "/api/v2.0/projects/" + p + "/repositories"
			req, err := http.NewRequest("GET", u, nil)
			if err != nil {
				panic(err)
			}
			params := make(url.Values)
			params.Set("page_size", "100")
			params.Set("page", s)
			req.URL.RawQuery = params.Encode()
			req.SetBasicAuth(i.User, i.Pass)
			r, err := http.DefaultClient.Do(req)
			body, err := ioutil.ReadAll(r.Body)

			var rr []map[string]interface{}
			json.Unmarshal(body, &rr)

			for _, v := range rr {
				b, _ := json.Marshal(v)
				mapStr, err := simplejson.NewJson(b)
				if err != nil {
					panic(err)
				}

				ime, err := mapStr.Get("name").String()
				if err != nil {
					panic(err)
				}
				arrayimage = append(arrayimage, ime)
			}
		}
	}
	return arrayimage, err
}

// 根据镜像信息，找到镜像中的 tag ，并写入至提供的 file 中
func (i *Infomation) HarborTag(ime []string, fileName string) error {
	var err error
	var imageTag ImageTag
	dstFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer dstFile.Close()
	log.Println("正在输出 harbor image tag 信息！！！\n")
	for _, v := range ime {
		u := i.Schema + "://" + i.URL + "/v2/" + v + "/tags/list"
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			panic(err)
		}
		params := make(url.Values)
		params.Set("page_size", "300")
		req.URL.RawQuery = params.Encode()
		req.SetBasicAuth(i.User, i.Pass)
		r, err := http.DefaultClient.Do(req)
		body, err := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &imageTag)
		for _, t := range imageTag.Tags {
			tag := i.URL + "/" + imageTag.Name + ":" + t
			fmt.Println(tag)
			dstFile.WriteString(tag + "\n")
		}
	}
	fmt.Println("")
	log.Printf("harbor image 信息输出至 %v 文件中，注意每次执行命令会覆盖之前的结果！！！", fileName)
	//return arraytag, err
	return err
}

// harbor cmd 的参数信息
func (i *Infomation) HarborCmd() {
	var repositry, file string
	flag.StringVar(&i.Schema, "schema", "http", "http or https，default http")
	flag.StringVar(&i.URL, "url", "harbor.k8s.local", "harbor Address，deault harbor.k8s.local")
	flag.StringVar(&i.User, "user", "admin", "harbor admin, default admin")
	flag.StringVar(&i.Pass, "passwd", "123456", "harbor password, default 123456")
	flag.StringVar(&repositry, "repositry", "all", "指定要选中的仓库，默认为 all，全部仓库")
	flag.StringVar(&file,"file", "harborImageList.txt", "运行结果输出文件")
	flag.Parse()

	if flag.NFlag() == 0 {
		log.Fatal("此命令必须传入参数, 否则无法执行\n  -passwd string\n        harbor password, default 123456 (default \"123456\")\n  -repositry string\n        指定要选中的仓库，默认为 all，全部仓库 (default \"all\")\n  -schema string\n        http or https，default http (default \"http\")\n  -url string\n        harbor Address，deault harbor.k8s.local (default \"harbor.k8s.local\")\n  -user string\n        harbor admin, default admin (default \"admin\")\n\n示例: harborImageTag --schema http --url harbor.k8s.local --user admin --passwd 123456 --repositry all --file harborImageList.txt\n")
	}

	if repositry == "all" {
		project, err := i.HarborProject()
		if err != nil {
			panic(err)
		}

		imagelist, err := i.HarborImage(project)
		if err != nil {
			panic(err)
		}

		i.HarborTag(imagelist, file)
	} else {
		imagelist, err := i.HarborImage([]string{repositry})
		if err != nil {
			panic(err)
		}

		i.HarborTag(imagelist, file)
	}

}

func main() {
	var i Infomation
	i.HarborCmd()
}
