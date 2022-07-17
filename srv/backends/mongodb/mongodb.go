package mongodb

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/orzen/steve/srv/plugin"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PluginMongo struct {
	plugin.BackendPlugin

	Cols     map[string]*mongo.Collection
	Client   *mongo.Client
	Database *mongo.Database
}

func New() plugin.BackendPlugin {
	return &PluginMongo{}
}

func (p *PluginMongo) Cfg() plugin.Cfg {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:     "mongo-host",
			Value:    "",
			Usage:    "MongoDB host",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "mongo-port",
			Value:    27017,
			Usage:    "MongoDB port",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "mongo-database",
			Value:    "",
			Usage:    "MongoDB database",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "mongo-user",
			Value:    "",
			Usage:    "MongoDB user",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "mongo-password",
			Value:    "",
			Usage:    "MongoDB password",
			Required: true,
		},
	}
	return plugin.Cfg{
		Name: "mongo",
		Cli:  flags,
	}
}

func (p *PluginMongo) Register() (plugin.Cfg, error) {
	log.Debug().Str("plugin", "mongo").Str("func", "register").Send()

	return p.Cfg(), nil
}

func (p *PluginMongo) Start(c *cli.Context) error {
	log.Debug().Str("plugin", "mongo").Str("func", "start").Send()

	user := c.String("mongo-user")
	pw := c.String("mongo-password")
	host := c.String("mongo-host")
	port := c.Int("mongo-port")
	dbName := c.String("mongo-database")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", user, pw, host, port)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return fmt.Errorf("establish database connection: %v", err)
	}

	db := client.Database(dbName)

	p.Client = client
	p.Database = db

	log.Info().
		Str("plugin", "mongo").
		Str("func", "start").
		Msgf("connected to database '%s'", uri)

	return nil
}

func (p *PluginMongo) Stop() error {
	log.Debug().Str("plugin", "mongo").Str("func", "stop").Send()

	if err := p.Client.Disconnect(context.TODO()); err != nil {
		log.Error().Err(err).Str("plugin", "mongo").Msg("disconnect")
		return err
	}

	return nil
}

type Field struct {
	Name   string
	Value  interface{}
	IsZero bool
}

func StructToMap(m interface{}) (map[string]Field, error) {
	ret := map[string]Field{}
	v := reflect.ValueOf(m)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ret, errors.New("value must be a struct")
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := typ.Field(i)
		fv := reflect.ValueOf(f)

		ret[f.Name] = Field{
			Name:   f.Name,
			Value:  fv.Interface(),
			IsZero: fv.IsZero(),
		}

	}

	return ret, nil
}

func MetaToFilter(m interface{}) (bson.D, error) {
	meta, err := StructToMap(m)
	if err != nil {
		return nil, errors.New("convert meta to map")
	}

	filter := bson.D{}
	for _, v := range meta {
		if !v.IsZero {
			e := primitive.E{
				Key:   v.Name,
				Value: v.Value,
			}
			filter = append(filter, e)
		}
	}

	return filter, nil
}

func (p *PluginMongo) Set(t string, r interface{}) error {
	log.Debug().Str("plugin", "mongo").Str("func", "set").Send()

	col := p.Database.Collection(t)

	if _, err := col.InsertOne(context.TODO(), r); err != nil {
		log.Error().Err(err).
			Str("op", "set").
			Str("type", t).
			Interface("resource", r).
			Msg("insert one")
		return plugin.ErrInternal
	}

	return nil
}

func (p *PluginMongo) Get(t string, m interface{}, retR interface{}) error {
	log.Debug().Str("plugin", "mongo").Str("func", "get").Send()

	col := p.Database.Collection(t)

	filter, err := MetaToFilter(m)
	if err != nil {
		log.Error().Err(err).
			Str("op", "get").
			Str("type", t).
			Interface("filter", m).
			Send()
		return plugin.ErrInternal
	}

	if err := col.FindOne(context.TODO(), filter).Decode(retR); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Error().Err(err).
				Str("op", "get").
				Str("type", t).
				Interface("filter", m).
				Msg("find one")
			return plugin.ErrInternal
		}
	}

	return nil
}

func (p *PluginMongo) List(t string, m interface{}, retM interface{}) error {
	log.Debug().Str("plugin", "mongo").Str("func", "list").Send()

	col := p.Database.Collection(t)

	filter, err := MetaToFilter(m)
	if err != nil {
		log.Error().Err(err).
			Str("op", "list").
			Str("type", t).
			Interface("filter", m).
			Send()
		return plugin.ErrInternal
	}

	cur, err := col.Find(context.TODO(), filter)
	if err != nil {
		log.Warn().Err(err).
			Str("op", "list").
			Str("type", t).
			Interface("filter", m).
			Msg("find document")
		return plugin.ErrInternal
	}
	defer cur.Close(context.TODO())

	if err = cur.All(context.TODO(), retM); err != nil {
		log.Warn().Err(err).
			Str("op", "list").
			Str("type", t).
			Interface("filter", m).
			Msg("cursor all")
		return plugin.ErrInternal
	}

	return nil
}

func (p *PluginMongo) Delete(t string, m interface{}, retR interface{}) error {
	log.Debug().Str("plugin", "mongo").Str("func", "delete").Send()

	col := p.Database.Collection(t)

	filter, err := MetaToFilter(m)
	if err != nil {
		log.Error().Err(err).
			Str("op", "list").
			Str("type", t).
			Interface("filter", m).
			Send()
		return plugin.ErrInternal
	}

	if err := col.FindOneAndDelete(context.TODO(), filter).Decode(&retR); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Error().Err(err).
				Str("op", "list").
				Str("type", t).
				Interface("filter", m).
				Msg("find one and delete")
			return plugin.ErrInternal
		}
	}

	return nil
}
