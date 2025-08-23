package categorymodel

import (
	"crud-go/config"
	"crud-go/entities"
)

func GetAll() []entities.Category {
	rows, err := config.DB.Query("SELECT * FROM categories")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	var categories []entities.Category
	
	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}

	return categories
}

func Create(category entities.Category) bool{
	result, err :=config.DB.Exec(`
		INSERT INTO categories (name, created_at, updated_at)
		VALUES (?, ?, ?)`,
		category.Name, category.CreatedAt, category.UpdatedAt,
	)

	if err != nil {
		panic(err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	return lastInsertId > 0
}

func GetById(id string) entities.Category {
	var category entities.Category
	err := config.DB.QueryRow("SELECT * FROM categories WHERE id = ?", id).Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		panic(err)
	}
	return category
}

func Update(category entities.Category, id string) bool {
	result, err := config.DB.Exec(`
		UPDATE categories
		SET name = ?, updated_at = ?
		WHERE id = ?`,
		category.Name, category.UpdatedAt, id,
	)

	if err != nil {
		panic(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return rowsAffected > 0
}