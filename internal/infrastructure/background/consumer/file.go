package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
)

func (c *Consumer) startFileConsumer() {
	go c.startDeleteFile()
}

func (c *Consumer) startDeleteFile() {
	if err := c.mqPro.ConsumeMessage(port.QueueNameDeleteFile, port.ExchangeFile, port.RoutingKeyDeleteFile, port.ExchangeTypeDirect, func(body []byte) error {
		var deleteFileMsg data.DeleteFileMessage
		if err := json.Unmarshal(body, &deleteFileMsg); err != nil {
			c.logPro.Error("Json unmarshal delete files message failed", port.String("schema", deleteFileMsg.Schema), port.Error(err))
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := c.storagePro.DeleteBatch(ctx, deleteFileMsg.Keys); err != nil {
			c.logPro.Error("Batch delete file failed", port.String("schema", deleteFileMsg.Schema), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		c.logPro.Error("Start consumer delete files failed", port.Error(err))
	}
}
