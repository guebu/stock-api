package elasticsearch

import(
	"context"
	"fmt"
	"github.com/guebu/common-utils/logger"
	"github.com/olivere/elastic"
	"time"
)

var(
	Client myEsClientInterface = &myEsClient{}
)

type myEsClientInterface interface {
	setClient(*elastic.Client)
	Index(string, string, interface{}) (*elastic.IndexResponse, error)
	Get(string, string, string) (*elastic.GetResult, error)
	Search(string, elastic.Query) (*elastic.SearchResult, error)
}

type myEsClient struct {
	client *elastic.Client
}

func Init() {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.0.1:9200"),
		elastic.SetSniff(false),
	    elastic.SetHealthcheckInterval(5*time.Second),
		)

	if err != nil {
		fmt.Println(err)
		panic( err )
	}

	Client.setClient(client)

}

func (c *myEsClient) setClient(client *elastic.Client){
	c.client = client
}

func (c *myEsClient) Index(index string, doctype string, doc interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()

	// Index the given object
	put1, err := c.client.Index().
		Index(index).
		Type(doctype).
		BodyJson(doc).
		Do(ctx)

	if err != nil {
		// Handle error
		logger.Error("Error when trying to index given document...", err, "Layer:ES-Client", "Status:Error")
		fmt.Println(err)
		return nil, err
	}
	logger.Info(fmt.Sprintf("Indexed stock item %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type), "App:Stock-Api", "Layer:ES-Client", "Func:Index", "Status:End")
	return put1, nil
}

func (c *myEsClient) Get(index string, doctype string, id string) (*elastic.GetResult, error) {
	// Get tweet with specified ID
	ctx := context.Background()
	result, err := c.client.Get().
		Index(index).
		Type(doctype).
		Id(id).
		Do(ctx)
	if err != nil {
		// Handle error
		logger.Error(fmt.Sprintf("Error when trying to get stock item from ES with id %s", id), err, "App:Stock-Api", "Layer:ES-Client", "Func:Get", "Status:Error")
		return nil, err
	}

	if result.Found {
		logger.Info( fmt.Sprintf("Got document %s in version %d from index %s, type %s\n", result.Id, result.Version, result.Index, result.Type), "App:Stock-Api", "Layer:ES-Client", "Func:Get", "Status:End")
	} else {
		logger.Info( fmt.Sprintf("No document found with given data -- ID: %s, Version: %d, Index:  %s, Type %s\n", result.Id, result.Version, result.Index, result.Type), "App:Stock-Api", "Layer:ES-Client", "Func:Get", "Status:End")
	}

	return result, nil
}

func (c *myEsClient) Search(index string, query elastic.Query) (*elastic.SearchResult, error) {

	logger.Info("Start searching stock items...", "App:StockItem-API", "Layer:ES-Client", "Func:Search", "Status:Start")
	ctx := context.Background()
	result, err := c.client.Search(index).Query(query).RestTotalHitsAsInt(true).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Error when trying to search documents in index %s.", index), err, "App:StockItem-API", "Layer:ES-Client", "Func:Search", "Status:Error")
		return nil, err
	}
	logger.Info("Searched successfully for stock itesm...", "App:StockItem-API", "Layer:ES-Client", "Func:Search", "Status:End")
	return result, nil

}

