package db

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"labix.org/v2/mgo/bson"

	"github.com/stretchr/testify/assert"
)

var (
	db_test       *DB
	brhub_test    *Brhub
	item_test     *Item
	comments_test []*Comment
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().UnixNano())

	session := MainSession()
	defer session.Close()

	db_test = session.DB()

	m.Run()
}

func TestAddBrhub(t *testing.T) {
	brhub_test = &Brhub{Id: bson.NewObjectId(), Name: "Test"}
	assert.Nil(t, db_test.AddBrhub(brhub_test))
}

func TestAddItem(t *testing.T) {
	item_test = NewItem()
	assert.Nil(t, db_test.AddItem(item_test))
}

func TestAddComment(t *testing.T) {
	comnb := 1 + rand.Intn(4)
	comments_test = make([]*Comment, comnb)

	for i := comnb - 1; i >= 0; i-- {
		comments_test[i] = NewComment()
		comments_test[i].Item = item_test.Id
		comments_test[i].Comments = []*Comment{}
		assert.Nil(t, db_test.AddComment(comments_test[i]))

		assert.Nil(t, saveTree(comments_test[i], i))
	}
}

func TestCommentTree(t *testing.T) {
	errc := make(chan error, 1)
	go db_test.commentTree(item_test, errc)
	assert.Nil(t, <-errc)

	assert.True(t, reflect.DeepEqual(item_test.Comments, comments_test))
}

func saveTree(c *Comment, depth int) error {
	var err error

	n := 1 + rand.Intn(4)
	c.Comments = make([]*Comment, n)
	for i := n - 1; i >= 0; i-- {
		child := NewComment()
		child.Parent = c.Id
		child.Item = c.Item
		child.Comments = []*Comment{}

		err := db_test.AddComment(child)
		if err != nil {
			return err
		}
		c.Comments[i] = child
	}

	if depth > 0 {
		for _, comment := range c.Comments {
			if err != nil {
				break
			}
			err = saveTree(comment, depth-1)
		}
	}
	return err
}

func printTree(coms []*Comment, acc int) {
	for _, child := range coms {
		fmt.Printf(strings.Repeat(" ", acc)+"%+v\n", child)
		printTree(child.Comments, acc+5)
	}
}
