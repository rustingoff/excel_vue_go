package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"github.com/rustingoff/excel_vue_go/internal/models"
	"log"
	"time"
)

const _UserIndex = "users"

type UserRepository interface {
	CreateUser(user models.User) error
	DeleteUser(id string) error
	GetUserById(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)

	Login(email, token string) error
}

type userRepository struct {
	elasticClient *elastic.Client
	redisClient   *redis.Client
}

func GetUserRepository(e *elastic.Client, r *redis.Client) UserRepository {
	return &userRepository{elasticClient: e, redisClient: r}
} //nolint:typechecking

func (repo *userRepository) CreateUser(user models.User) error {
	_, err := repo.elasticClient.Index().Index(_UserIndex).BodyJson(user).Do(context.TODO())
	if err != nil {
		log.Println("[ERR]: failed to create user, ", err.Error())
		return err
	}

	return nil
}

func (repo *userRepository) DeleteUser(id string) error {

	_, err := repo.elasticClient.Delete().Index(_UserIndex).Id(id).Do(context.TODO())
	if err != nil {
		log.Println("[ERR]: failed to delete user, ", err.Error())
		return err
	}

	return nil
}

func (repo *userRepository) GetUserById(id string) (models.User, error) {
	query := elastic.NewMatchQuery("_id", id)

	res, err := repo.elasticClient.Search(_UserIndex).Query(query).Do(context.TODO())
	if err != nil {
		log.Println("[ERR]: failed to get user by id")
		return models.User{}, err
	}

	if int(res.TotalHits()) > 0 {
		var user models.User

		err = json.Unmarshal(res.Hits.Hits[0].Source, &user)
		if err != nil {
			log.Println("[ERR]: failed to unmarshal source")
			return models.User{}, err
		}
		user.ID = res.Hits.Hits[0].Id
		return user, nil
	}

	return models.User{}, errors.New("user not found")
}

func (repo *userRepository) GetUserByEmail(email string) (models.User, error) {
	query := elastic.NewMatchQuery("email", email)

	res, err := repo.elasticClient.Search(_UserIndex).Query(query).Do(context.TODO())
	if err != nil {
		log.Println("[ERR]: failed to get user by email")
		return models.User{}, err
	}

	if int(res.TotalHits()) > 0 {
		var user models.User

		err = json.Unmarshal(res.Hits.Hits[0].Source, &user)
		if err != nil {
			log.Println("[ERR]: failed to unmarshal source")
			return models.User{}, err
		}
		user.ID = res.Hits.Hits[0].Id
		return user, nil
	}

	return models.User{}, errors.New("user not found")
}

func (repo *userRepository) Login(email, token string) error {
	status := repo.redisClient.Set(context.TODO(), email, token, time.Hour*24)
	if status.Err() != nil {
		log.Println("[ERR]: failed to set value in redis")
		return status.Err()
	}

	return nil
}
