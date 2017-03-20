package mangafox

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/Zignd/bewitched-marmot/types"
	"github.com/pkg/errors"
)

const host = "www.mangafox.me"

var baseURL = fmt.Sprintf("http://%s", host)

// Search queries the site
func Search(query string) ([]*types.CompactManga, error) {
	// searchURL := fmt.Sprintf("%s/search.php?name_method=cw&name=

	// HERE HERE HERE

	// &type=&author_method=cw&author=&artist_method=cw&artist=&genres%5BAction%5D=0&genres%5BAdult%5D=0&genres%5BAdventure%5D=0&genres%5BComedy%5D=0&genres%5BDoujinshi%5D=0&genres%5BDrama%5D=0&genres%5BEcchi%5D=0&genres%5BFantasy%5D=0&genres%5BGender+Bender%5D=0&genres%5BHarem%5D=0&genres%5BHistorical%5D=0&genres%5BHorror%5D=0&genres%5BJosei%5D=0&genres%5BMartial+Arts%5D=0&genres%5BMature%5D=0&genres%5BMecha%5D=0&genres%5BMystery%5D=0&genres%5BOne+Shot%5D=0&genres%5BPsychological%5D=0&genres%5BRomance%5D=0&genres%5BSchool+Life%5D=0&genres%5BSci-fi%5D=0&genres%5BSeinen%5D=0&genres%5BShoujo%5D=0&genres%5BShoujo+Ai%5D=0&genres%5BShounen%5D=0&genres%5BShounen+Ai%5D=0&genres%5BSlice+of+Life%5D=0&genres%5BSmut%5D=0&genres%5BSports%5D=0&genres%5BSupernatural%5D=0&genres%5BTragedy%5D=0&genres%5BWebtoons%5D=0&genres%5BYaoi%5D=0&genres%5BYuri%5D=0&released_method=eq&released=&rating_method=eq&rating=&is_completed=&advopts=1", baseURL, url.QueryEscape(query))

	// req, err := http.NewRequest("GET", searchURL, nil)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "Search(\"%s\") could not create request", searchURL)
	// }

	// req.Header.Set("Referer", baseURL)

	// res, err := http.DefaultClient.Do(req)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "Search(\"%s\") failed to perform query", searchURL)
	// }
	// defer res.Body.Close()

	// doc, err := goquery.NewDocumentFromResponse(res)
	// if err != nil {
	// 	return nil, errors.Wrapf(err, "Search(\"%s\") failed to parse HTML document", searchURL)
	// }
	// if wasThrottled(doc) == true {
	// 	duration := time.Duration(5) * time.Second
	// 	time.Sleep(duration)
	// 	return Search(query)
	// }

	// result := []*types.CompactManga{}

	// docResultItems := doc.Find()

	return nil, errors.New("Not implemented")
}

func wasThrottled(doc *goquery.Document) bool {
	for _, node := range doc.Find("#page > div.left > div.border.clear").Nodes {
		if node.Data == "Sorry, can‘t search again within 5 seconds." {
			return true
		}
	}
	return false
}

// GetManga retrieves a mangafox
func GetManga(mangaURL string) (*types.DetailedManga, error) {
	return nil, errors.New("Not implemented")
}

// GetChapter retrieves a GetChapter
func GetChapter(chapterURL string) (*types.DetailedChapter, error) {
	return nil, errors.New("Not implemented")
}
