package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/ercole-io/tico/service"
)

func OracleCloudHandler(ctx context.Context, in io.Reader, out io.Writer) {
	res := service.ExampleSearchResources()
	for _, v := range res {
		json.NewEncoder(out).Encode(fmt.Sprintf("Identifier:%v, DisplayName:%v, DefinedTags:%v", v.Identifier, v.DisplayName, v.DefinedTags))
	}

}
