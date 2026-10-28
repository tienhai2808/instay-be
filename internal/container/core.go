package container

import (
	"fmt"

	"github.com/InstaySystem/is_v2-be/internal/infrastructure/initialization"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/ai"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/apikey"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/didit"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/fcm"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/gcs"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/jwt"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/orm"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/rabbitmq"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/redis"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/smtp"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/sonyflake"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/zap"
)

func (c *Container) initCore() (err error) {
	c.log, err = initialization.InitLog(c.cfg.Log)
	if err != nil {
		return err
	}

	c.SQL, err = initialization.InitSQL(c.cfg.SQL)
	if err != nil {
		return err
	}

	c.NoSQL, err = initialization.InitNoSQL(c.cfg.NoSQL)
	if err != nil {
		return err
	}

	c.cache, err = initialization.InitCache(c.cfg.Cache)
	if err != nil {
		return err
	}

	c.stor, err = initialization.InitStorage(c.cfg.Storage)
	if err != nil {
		return err
	}

	c.mq, err = initialization.InitMessageQueue(c.cfg.MQ)
	if err != nil {
		return err
	}

	c.pn, err = initialization.InitPushNotification()
	if err != nil {
		return err
	}

	c.idGen, err = initialization.InitIDGen()
	if err != nil {
		return err
	}

	aiAddr := fmt.Sprintf("%s:%d", c.cfg.Services.AIHost, c.cfg.Services.AIPort)

	ca := initialization.ClientAddresses{
		AIAddr: aiAddr,
	}

	c.gRPCClients, err = initialization.InitClients(ca)
	if err != nil {
		return err
	}

	c.LogPro = zap.NewLogProvider(c.log.Logger())

	c.IDPro = sonyflake.NewIDProvider(c.idGen.Generator())

	c.tokenPro = jwt.NewTokenProvider(c.cfg.Token)
	c.apiKeyPro = apikey.NewAPIKeyProvider(c.cfg.Token)

	c.CachePro = redis.NewCacheProvider(c.cache.Client())

	c.MQPro = rabbitmq.NewMessageQueueProvider(c.mq.Connection(), c.log.Logger())

	c.MailerPro = smtp.NewMailerProvider(c.cfg.Mailer)

	c.txPro = orm.NewTxProvider(c.SQL.ORM())

	c.StoragePro = gcs.NewStorageProvider(c.cfg.Storage, c.stor.Client())

	c.PNPro = fcm.NewPushNotificationProvider(c.pn.Client())

	c.KYCPro = didit.NewKYCProvider(c.cfg.KYC)

	c.AIPro = ai.NewAIProvider(c.gRPCClients.AIClient)

	return nil
}
