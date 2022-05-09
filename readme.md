# Table of Content
<ol>
    <li>
        <a href="">Name</a>
    </li>
    <li>
        <a href="">Getting started</a>
        <ul>
            <a href="">Installation</a>
        </ul>
    </li>
    <li>
        <a href="">Usage</a>
    </li>
</ol>

# SQLABST
sqlsqbdr is the acronym for SQL Simple Query Builder. This is a simple query builder that only build the placeholder for updating field, inserting field, and filtering field. No more creating a long placeholder of query for inserting, and updating data 


# Getting started
This is an example to how to use this project locally

## Installation
    go get github.com/nurcahyaari/sqlsqbdr

# Motivasion

when I created a query that have a long placeholder for inserting and updating data. it's painful, because there's many field, placeholder, and the value that need to assign to the query. especially when I need to add new field (it means new field from DB, or missing field from query). I feel confused by the query

here the example when I create query
```go

INSERT INTO product 
(
    id, 
    name,
    image, 
    price, 
    discount_price, 
    active, 
    created_at, 
    created_by, 
    updated_at, 
    updated_by
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
)
```
So, I want to create the query without defining the field, and their placeholder


# Usage

This is a guideline to use sqlsqbdr on your projects

## example
```go
// let's say you have a defined struct
type Product struct {
    Id int64 `db:"id"`
    Name string `db:"name"`
    Image string `db:"image"`
    Price int64 `db:"price"`
    DiscountPrice int64 `db:"discount_price"`
    Active bool `db:"active"`
    CreatedAt string `db:"created_at"`
    CreatedBy string `db:"created_by"`
    UpdatedAt string `db:"updated_at"`
    UpdatedBy string `db:"updated_by"`
}
// here is my new data
products := []Product{
    Product{
        Id: 1,
        Name: "Test",
        Image: "test.com/1.png",
        Price: 1000,
        DiscountPrice: 1000,
        Active: true,
        CreatedAt: "2022-01-01 01:00:00"
        CreatedBy: "Tester",
        UpdatedAt: "2022-01-01 01:00:00"
        UpdatedBy: "Tester",
    }
}

// here is my updated data
updatedProduct := Product{
    Name: "Test2",
    Price: 2000,
    DiscountPrice: 1500,
}
```

### Type of FieldSelect
sqlsqbdr have 2 field type
<table border="1">
<tr>
    <th>name</th>
    <th>usage</th>
</tr>
<tr>
    <td>IncludeField</td>
    <td>Select field based on the selected field on the params</td>
</tr>
<tr>
    <td>ExcludeField</td>
    <td>Select field that not in selected field on the params</td>
</tr>
</table>

### build insert placeholder

`BuildInsertField` have 3 params first is the model/ entity's struct, the second is FieldSelectType, the the last is a variadic this params is defined the list of field that you want to insert if you choose the FieldSelectType as `IncludeField`. but if you chose the FieldSelectedType as `ExcludeField` the last params is the field that you don't want to insert



```go
insertField, err := sqlsqbdr.BuildInsertField(products, sqlsqbdr.IncludeField)
if err != nil {
    return
}

query := fmt.Sprintf("INSERT INTO product (%s) VALUES %s", strings.Join(insertField.Name, ","), strings.Join(insertField.Placeholder, ","))

// the query will be like this
/*
    INSERT INTO product 
    (
        id, 
        name,
        image, 
        price, 
        discount_price, 
        active, 
        created_at, 
        created_by, 
        updated_at, 
        updated_by
    ) VALUES (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
    )
*/
```


### build update placeholder
`BuildUpdatedField` have 3 params. first is the model/ entity's struct, the second is FieldSelectType, the the last is a variadic this params is defined the list of field that you want to update if you choose the FieldSelectType as `IncludeField`. but if you chose the FieldSelectedType as `ExcludeField` the last params is the field that you don't want to update

```go
updatedField, err := sqlsqbdr.BuildUpdatedField(updatedProduct, sqlsqbdr.IncludeField, "name", "price", "discount_price")
if err != nil {
    return
}

query := fmt.Sprintf("UPDATE product SET %s", strings.Join(updatedField.Name, ","))
```

### build filtering field
`BuildWhereFilter` has 1 param. and it's a defined struct

#### the parameter
`Filters` is an array, and here is the description about the struct

<table border="1">
<tr>
    <th>name</th>
    <th>usage</th>
</tr>
<tr>
    <td>Field</td>
    <td>it means the field on your database</td>
</tr>
<tr>
    <td>Value</td>
    <td>It means the value of the field</td>
</tr>
</table>

```go
field, values := sqlsqbdr.BuildWhereFilter(sqlsqbdr.Filters{
    &sqlsqbdr.Filter{
        Field: "id",
        Value: 1,
    },
})
// result
// field: id = ?
// values: []interface{}{1}
```