package mongo_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCountUsers(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	// defer mt.Stop()

	mt.Run("Count users", func(mt2 *mtest.T) {
		mt2.AddMockResponses(
			mtest.CreateCursorResponse(1, "test.users", mtest.FirstBatch, bson.D{{Key: "n", Value: int32(5)}}),
		)

		coll := mt2.Coll
		count, err := coll.CountDocuments(context.Background(), bson.M{})
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})
}
