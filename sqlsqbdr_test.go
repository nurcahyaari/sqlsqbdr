package sqlsqbdr_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nurcahyaari/sqlsqbdr"
	"github.com/stretchr/testify/assert"
)

// let's say you have a defined struct
type Product struct {
	Id            int64  `db:"id"`
	Name          string `db:"name"`
	Image         string `db:"image"`
	Price         int64  `db:"price"`
	DiscountPrice int64  `db:"discount_price"`
	Active        bool   `db:"active"`
	CreatedAt     string `db:"created_at"`
	CreatedBy     string `db:"created_by"`
	UpdatedAt     string `db:"updated_at"`
	UpdatedBy     string `db:"updated_by"`
}

func TestFullQueryCreational(t *testing.T) {
	t.Run("Create full query insert", func(t *testing.T) {

		// here is my new data
		products := []Product{
			{
				Id:            1,
				Name:          "Test",
				Image:         "test.com/1.png",
				Price:         1000,
				DiscountPrice: 1000,
				Active:        true,
				CreatedAt:     "2022-01-01 01:00:00",
				CreatedBy:     "Tester",
				UpdatedAt:     "2022-01-01 01:00:00",
				UpdatedBy:     "Tester",
			},
		}

		insertField, err := sqlsqbdr.BuildInsertField(products, sqlsqbdr.IncludeField)
		if err != nil {
			return
		}
		insertFieldValue := [][]interface{}{
			{
				int64(1),
				"Test",
				"test.com/1.png",
				int64(1000),
				int64(1000),
				true,
				"2022-01-01 01:00:00",
				"Tester",
				"2022-01-01 01:00:00",
				"Tester",
			},
		}

		query := fmt.Sprintf("INSERT INTO product (%s) VALUES %s", strings.Join(insertField.Name, ","), strings.Join(insertField.Placeholder, ","))
		queryExpected := "INSERT INTO product (id,name,image,price,discount_price,active,created_at,created_by,updated_at,updated_by) VALUES (?,?,?,?,?,?,?,?,?,?)"

		assert.Equal(t, insertField.Values, insertFieldValue)
		assert.Equal(t, queryExpected, query)
	})

	t.Run("Create full query update with filter", func(t *testing.T) {
		// here is my updated data
		updatedProduct := Product{
			Name:          "Test2",
			Price:         2000,
			DiscountPrice: 1500,
		}

		updatedField, err := sqlsqbdr.BuildUpdatedField(updatedProduct, sqlsqbdr.IncludeField, "name", "price", "discount_price")
		if err != nil {
			return
		}
		updatedFieldValueExp := []interface{}{
			"Test2",
			int64(2000),
			int64(1500),
		}

		whereField, whereValue := sqlsqbdr.BuildWhereFilter(sqlsqbdr.Filters{
			&sqlsqbdr.Filter{
				Field: "id",
				Value: updatedProduct.Id,
			},
		})
		whereValueExp := []interface{}{updatedProduct.Id}

		query := fmt.Sprintf("UPDATE product SET %s WHERE %s", strings.Join(updatedField.Name, ","), whereField)
		queryExpected := "UPDATE product SET name = ?,price = ?,discount_price = ? WHERE id = ?"
		assert.Equal(t, queryExpected, query)
		assert.Equal(t, updatedFieldValueExp, updatedField.Value)
		assert.Equal(t, whereValueExp, whereValue)
	})
}
