package experiment

import (
	"log"
	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
	"github.com/Obaied/RAKE.go"
	"github.com/tidinio/src/component/helper/string"
	"github.com/tidinio/src/component/helper/http"
)

func ItemKeywords() {

	//TODO: also get meta keywords (for example with GoOse)
	log.Println("tut: PredictItemTags a")
	url := "https://futurism.com/scientists-have-captured-the-first-ever-image-of-a-dark-matter-web/"
	htmlText := string_helper.StripHtmlContent(http_helper.GetContentFromUrl(url))
	log.Println("tut: PredictItemTags b")

	candidates := rake.RunRake(htmlText)
	popular := ""
	for _, candidate := range candidates {
		if (candidate.Value > 7) {
			popular = popular + " " +candidate.Key
		}
	}
	log.Println("tut: candidates ", popular)

	stream := make(chan base.TextDatapoint, 200)
	model := text.NewNaiveBayes(stream, 2, base.OnlyWordsAndNumbers)
	model.String()
	tf := text.TFIDF(*model)
	freq := tf.MostImportantWords(popular, 10)
	log.Println("tut: MostImportantWords ", freq)
}
