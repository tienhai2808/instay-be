package implement

import (
	"context"
	"errors"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mongo/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type messageRepositoryImpl struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) repository.MessageRepository {
	collection := db.Collection(entity.MessageCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
		{
			Keys: bson.D{
				{Key: "schema", Value: 1},
				{Key: "sender_type", Value: 1},
				{Key: "is_deleted", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "schema", Value: 1},
				{Key: "chat_id", Value: 1},
				{Key: "sender_type", Value: 1},
				{Key: "booking_room_read_at", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "schema", Value: 1},
				{Key: "chat_id", Value: 1},
				{Key: "sender_type", Value: 1},
				{Key: "message_users.user_id", Value: 1},
			},
		},
	}

	if _, err := collection.Indexes().CreateMany(ctx, indexes); err != nil {
		return &messageRepositoryImpl{collection}
	}

	return &messageRepositoryImpl{collection}
}

func (r *messageRepositoryImpl) Create(ctx context.Context, message *model.Message) error {
	msgEntity := mapper.ToMessageEntity(message)
	_, err := r.collection.InsertOne(ctx, msgEntity)
	return err
}

func (r *messageRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Message, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, exceptions.ErrInvalidID
	}

	var msg entity.MessageEntity
	if err := r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&msg); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToMessage(&msg), nil
}

func (r *messageRepositoryImpl) FindAllByChatIDPaginated(ctx context.Context, schema string, chatID int64, page, limit uint32) ([]*model.Message, int64, error) {
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	filter := bson.M{
		"schema":  schema,
		"chat_id": chatID,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return make([]*model.Message, 0), 0, nil
	}

	offset := int64(page-1) * int64(limit)
	findOpts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1},
		}).
		SetSkip(offset).
		SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var msgEntities []*entity.MessageEntity
	if err = cursor.All(ctx, &msgEntities); err != nil {
		return nil, 0, err
	}

	return mapper.ToMessages(msgEntities), total, nil
}

func (r *messageRepositoryImpl) FindAllByChatIDSinceCursor(ctx context.Context, schema string, chatID int64, since, cursorAt *time.Time, cursorID string, limit uint32) ([]*model.Message, bool, error) {
	if limit == 0 {
		limit = 100
	}
	if limit > 100 {
		limit = 100
	}

	andFilters := bson.A{
		bson.M{"schema": schema},
		bson.M{"chat_id": chatID},
	}

	if since != nil {
		andFilters = append(andFilters, bson.M{"created_at": bson.M{"$gte": *since}})
	}

	if cursorAt != nil {
		orFilters := bson.A{
			bson.M{"created_at": bson.M{"$lt": *cursorAt}},
		}

		if cursorID != "" {
			objID, err := bson.ObjectIDFromHex(cursorID)
			if err != nil {
				return nil, false, exceptions.ErrInvalidID
			}

			orFilters = append(orFilters, bson.M{
				"created_at": *cursorAt,
				"_id":        bson.M{"$lt": objID},
			})
		}

		andFilters = append(andFilters, bson.M{"$or": orFilters})
	}

	filter := bson.M{"$and": andFilters}
	findLimit := int64(limit) + 1

	findOpts := options.Find().
		SetSort(bson.D{
			{Key: "created_at", Value: -1},
			{Key: "_id", Value: -1},
		}).
		SetLimit(findLimit)

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, false, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var msgEntities []*entity.MessageEntity
	if err = cursor.All(ctx, &msgEntities); err != nil {
		return nil, false, err
	}

	hasNext := len(msgEntities) > int(limit)
	if hasNext {
		msgEntities = msgEntities[:int(limit)]
	}

	return mapper.ToMessages(msgEntities), hasNext, nil
}

func (r *messageRepositoryImpl) SetMessageReactionsByID(ctx context.Context, schema, id string, reactions []*model.MessageReaction) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return exceptions.ErrInvalidID
	}

	filter := bson.M{
		"_id":    objID,
		"schema": schema,
	}

	update := bson.M{
		"$set": bson.M{
			"message_reactions": mapper.ToMessageReactionsEntity(reactions),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return exceptions.ErrMessageNotFound
	}

	return nil
}

func (r *messageRepositoryImpl) SetBookingRoomReadAtAllByChatID(ctx context.Context, schema string, chatID int64, readAt time.Time) error {
	filter := bson.M{
		"schema":      schema,
		"chat_id":     chatID,
		"sender_type": model.MessageActionTypeGuest,
		"$or": bson.A{
			bson.M{"booking_room_read_at": bson.M{"$exists": false}},
			bson.M{"booking_room_read_at": nil},
		},
	}

	update := bson.M{
		"$set": bson.M{
			"booking_room_read_at": readAt,
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *messageRepositoryImpl) PushMessageUserAllByChatID(ctx context.Context, schema string, chatID int64, messageUser *model.MessageUser) error {
	filter := bson.M{
		"schema":                schema,
		"chat_id":               chatID,
		"sender_type":           model.MessageActionTypeGuest,
		"message_users.user_id": bson.M{"$ne": messageUser.UserID},
	}

	update := bson.M{
		"$push": bson.M{
			"message_users": mapper.ToMessageUserEntity(messageUser),
		},
	}

	_, err := r.collection.UpdateMany(ctx, filter, update)
	return err
}

func (r *messageRepositoryImpl) CountDashboardUnreadGuestMessages(ctx context.Context, schema string, chatIDs ...[]int64) (int64, error) {
	filter := bson.M{
		"schema":      schema,
		"sender_type": model.MessageActionTypeGuest,
		"is_deleted":  false,
		"$or": bson.A{
			bson.M{"message_users": bson.M{"$exists": false}},
			bson.M{"message_users": nil},
			bson.M{"message_users": bson.A{}},
		},
	}

	if len(chatIDs) > 0 {
		if len(chatIDs[0]) == 0 {
			return 0, nil
		}
		filter["chat_id"] = bson.M{"$in": chatIDs[0]}
	}

	return r.collection.CountDocuments(ctx, filter)
}
