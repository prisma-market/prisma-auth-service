package mongodb

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kihyun1998/prisma-market/prisma-auth-service/internal/models"
)

type AuthRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

// NewAuthRepository AuthRepository 생성자
func NewAuthRepository(mongoURI string) (*AuthRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB 연결
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	// 연결 테스트
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	db := client.Database("prisma_market")
	collection := db.Collection("auth_data")

	// 이메일 unique 인덱스 생성
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return &AuthRepository{
		db:         db,
		collection: collection,
	}, nil
}

// CreateUser 새로운 사용자 생성
func (r *AuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Status = "active"

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("email already exists")
		}
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return nil
}

// FindUserByEmail 이메일로 사용자 찾기
func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin 마지막 로그인 시간 업데이트
func (r *AuthRepository) UpdateLastLogin(ctx context.Context, userID primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"last_login": now,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateByID(ctx, userID, update)
	return err
}

// Close MongoDB 연결 종료
func (r *AuthRepository) Close(ctx context.Context) error {
	return r.db.Client().Disconnect(ctx)
}
