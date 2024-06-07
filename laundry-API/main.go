package main

import (
	"database/sql"
	"laundry/config"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber int    `json:"phoneNumber"`
	Address     string `json:"address"`
}

type Products struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Unit  string `json:"unit"`
}

type Employee struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

type Transaction struct {
	Id          string       `json:"id"`
	BillDate    string       `json:"billDate"`
	EntryDate   string       `json:"entryDate"`
	FinishDate  string       `json:"finishDate"`
	EmployeeID  string       `json:"employeeId"`
	CustomerID  string       `json:"customerId"`
	Employee    Employee     `json:"employee"`
	Customer    Customer     `json:"customer"`
	BillDetails []BillDetail `json:"billDetails"`
	TotalBill   int          `json:"totalBill"`
}

type BillDetail struct {
	ID            int      `json:"id"`
	TransactionID string   `json:"billId"`
	ProductID     string   `json:"productId"`
	Qty           int      `json:"qty"`
	Product       Products `json:"product"`
}

func main() {
	r := gin.Default()
	var db = config.ConnectDB()
	if db == nil {
		panic("Tidak dapat terhubung dengan database")
	}
	defer db.Close()

	v1 := r.Group("/api")
	{
		customerGroup := v1.Group("/customer")
		{
			customerGroup.GET("/", getCustomers(db))
			customerGroup.GET("/:id", getCustomerById(db))
			customerGroup.POST("/", createCustomer(db))
			customerGroup.PUT("/:id", updateCustomer(db))
			customerGroup.DELETE("/:id", deleteCustomer(db))
		}

		productsGroup := v1.Group("/products")
		{
			productsGroup.GET("/", getProduct(db))
			productsGroup.GET("/:id", getProductById(db))
			productsGroup.POST("/", createProduct(db))
			productsGroup.PUT("/:id", updateProduct(db))
			productsGroup.DELETE("/:id", deleteProduct(db))
		}
		transactionGroup := v1.Group("/transaction")
		{
			transactionGroup.POST("/", createTransaction(db))
			transactionGroup.GET("/", getAllTransaction(db))
			transactionGroup.GET("/:id", getTransactionById(db))
		}
	}

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}

}

// get customerById
func getCustomerById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ngambil berdasarkan id
		id := c.Param("id")
		// mengconvert string ke id
		customerID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID customer tidak valid"})
			return
		}
		// mengambil struct dari customer
		var customer Customer
		// untuk melakukan koneksi database dengean perintah select
		row := db.QueryRow("SELECT no_hp, nama_customer, address FROM tbl_customer WHERE id = $1", customerID)
		// untuk melakukan validasi table
		err = row.Scan(&customer.PhoneNumber, &customer.Name, &customer.Address)
		// mengecek error
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data customer"})
			// untuk keluar dari function
			return
		}
		// mengambil ID dan mengubahnya ke dalam id yg sudah di convert
		customer.Id = customerID
		// jika ya maka tampilkan customer
		c.JSON(http.StatusOK, customer)
	}
}

// get customer
func getCustomers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// mengambil seluruh data customer
		var customers []Customer
		// melakukan pemanggilan database dengan query
		rows, err := db.Query("SELECT * FROM tbl_customer")
		// melakukan pengecekan jika error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			// untuk keluar dari func
			return
		}
		// untuk keluar dari database setelah semua di eksekusi
		defer rows.Close()
		// setelah itu melakukan otorisasi dengan menjalankan rows
		for rows.Next() {
			var customer Customer
			err := rows.Scan(&customer.Id, &customer.Name, &customer.PhoneNumber, &customer.Address)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan rows"})
				return
			}
			customers = append(customers, customer)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate rows"})
			return
		}

		c.JSON(http.StatusOK, customers)
	}
}

// create customer
func createCustomer(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var customer Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if customer.PhoneNumber <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "NoHp harus lebih dari 0"})
			return
		}

		if len(customer.Name) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "NamaCustomer harus memiliki minimal 3 karakter"})
			return
		}

		_, err := db.Exec("INSERT INTO tbl_customer ( no_hp, nama_customer, address) VALUES ($1, $2, $3)",
			customer.PhoneNumber, customer.Name, customer.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Customer berhasil ditambahkan"})
	}
}

// update customer
func updateCustomer(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		noHP := c.Param("id")
		customerId, err := strconv.Atoi(noHP)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer Id"})
			return
		}
		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM tbl_customer WHERE id = $1", customerId)
		err = row.Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus customer"})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
			return
		}

		var customer Customer
		if err := c.ShouldBindJSON(&customer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if customer.PhoneNumber <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "NoHp harus lebih dari 0"})
			return
		}
		if len(customer.Name) < 3 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "NamaCustomer harus memiliki minimal 3 karakter"})
			return
		}
		_, err = db.Exec("UPDATE tbl_customer SET nama_customer = $1, no_hp = $2, address = $3 WHERE id = $4", customer.Name, customer.PhoneNumber, customer.Address, customerId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil diupdate"})
	}
}

// Delete Customer
func deleteCustomer(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		customerID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
			return
		}
		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM tbl_customer WHERE id = $1", customerID)
		err = row.Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus customer"})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
			return
		}
		_, err = db.Exec("DELETE FROM tbl_customer WHERE id = $1", customerID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil dihapus"})
	}
}

func getProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []Products
		rows, err := db.Query("SELECT * FROM tbl_Products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var product Products
			err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Unit)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan rows"})
				return
			}
			products = append(products, product)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate rows"})
			return
		}

		c.JSON(http.StatusOK, products)
	}
}

// get customerById
func getProductById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		productID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID products tidak valid"})
			return
		}
		var product Products
		row := db.QueryRow("SELECT id, name, price, unit FROM tbl_products WHERE id = $1", productID)
		err = row.Scan(&product.ID, &product.Name, &product.Price, &product.Unit)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "products tidak ditemukan"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data products"})
			return
		}

		product.ID = productID
		c.JSON(http.StatusOK, product)
	}
}
func createProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product Products
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if product.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Price harus lebih dari 0"})
			return
		}

		if len(product.Name) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product harus memiliki minimal 5 karakter"})
			return
		}

		_, err := db.Exec("INSERT INTO tbl_products (name, price, unit) VALUES ($1, $2, $3)", product.Name, product.Price, product.Unit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan produk"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Produk berhasil ditambahkan"})
	}
}
func updateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		productID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID produk tidak valid"})
			return
		}

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM tbl_products WHERE id = $1", productID)
		err = row.Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate produk"})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}

		var product Products
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if product.Price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Harga harus lebih dari 0"})
			return
		}
		if len(product.Name) < 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nama produk harus memiliki minimal 5 karakter"})
			return
		}

		_, err = db.Exec("UPDATE tbl_products SET name = $1, price = $2, unit = $3 WHERE id = $4", product.Name, product.Price, product.Unit, productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate produk"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diupdate"})
	}
}

// Delete Product
func deleteProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		productID, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}
		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM tbl_products WHERE id = $1", productID)
		err = row.Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}
		_, err = db.Exec("DELETE FROM tbl_products WHERE id = $1", productID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
	}
}

// create transaction
func createTransaction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTransaction struct {
			BillDate    string       `json:"billDate"`
			EntryDate   string       `json:"entryDate"`
			FinishDate  string       `json:"finishDate"`
			EmployeeID  string       `json:"employeeId"`
			CustomerID  string       `json:"customerId"`
			BillDetails []BillDetail `json:"billDetails"`
		}

		if err := c.ShouldBindJSON(&newTransaction); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		employeeID, err := strconv.Atoi(newTransaction.EmployeeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employeeId"})
			return
		}

		customerID, err := strconv.Atoi(newTransaction.CustomerID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customerId"})
			return
		}
		if employeeID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Employee ID is required"})
			return
		}
		if customerID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
			return
		}
		if len(newTransaction.BillDetails) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No bill details provided"})
			return
		}
		var totalBill int
		for _, detail := range newTransaction.BillDetails {
			var price int
			err := db.QueryRow("SELECT price FROM tbl_products WHERE id = $1", detail.ProductID).Scan(&price)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			totalBill += price * detail.Qty
		}
		query := `
		INSERT INTO transactions (bill_date, entry_date, finish_date, employee_id, customer_id, total_bill)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
		var transactionID int
		err = db.QueryRow(query, newTransaction.BillDate, newTransaction.EntryDate, newTransaction.FinishDate, employeeID, customerID, totalBill).Scan(&transactionID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, detail := range newTransaction.BillDetails {
			_, err := db.Exec(`
			INSERT INTO bill_details (transaction_id, product_id, qty)
			VALUES ($1, $2, $3)
		`, transactionID, detail.ProductID, detail.Qty)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		c.JSON(http.StatusCreated, gin.H{
			"message":       "Transaction created successfully",
			"transactionID": transactionID,
			"totalBill":     totalBill,
		})
	}
}
func getTransactionById(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		row := db.QueryRow("SELECT * FROM transactions WHERE id = $1", id)
		transaction := Transaction{}
		err := row.Scan(
			&transaction.Id,
			&transaction.BillDate,
			&transaction.EntryDate,
			&transaction.FinishDate,
			&transaction.EmployeeID,
			&transaction.CustomerID,
			&transaction.TotalBill,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		employee := Employee{}
		err = db.QueryRow("SELECT * FROM tbl_employees WHERE id = $1", transaction.EmployeeID).Scan(
			&employee.ID,
			&employee.Name,
			&employee.PhoneNumber,
			&employee.Address,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		customer := Customer{}
		err = db.QueryRow("SELECT * FROM tbl_customer WHERE id = $1", transaction.CustomerID).Scan(
			&customer.Id,
			&customer.Name,
			&customer.PhoneNumber,
			&customer.Address,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := db.Query("SELECT * FROM bill_details WHERE transaction_id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		billDetails := []BillDetail{}
		for rows.Next() {
			billDetail := BillDetail{}
			err := rows.Scan(
				&billDetail.ID,
				&billDetail.TransactionID,
				&billDetail.ProductID,
				&billDetail.Qty,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			product := Products{}
			err = db.QueryRow("SELECT * FROM tbl_products WHERE id = $1", billDetail.ProductID).Scan(
				&product.ID,
				&product.Name,
				&product.Price,
				&product.Unit,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			billDetail.Product = product
			billDetails = append(billDetails, billDetail)
		}
		transaction.Employee = employee
		transaction.Customer = customer
		transaction.BillDetails = billDetails

		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction detail",
			"data":    []Transaction{transaction},
		})
	}
}
func getAllTransaction(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM transactions")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		transactions := []Transaction{}
		for rows.Next() {
			transaction := Transaction{}
			err := rows.Scan(
				&transaction.Id,
				&transaction.BillDate,
				&transaction.EntryDate,
				&transaction.FinishDate,
				&transaction.EmployeeID,
				&transaction.CustomerID,
				&transaction.TotalBill,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			employee := Employee{}
			err = db.QueryRow("SELECT * FROM tbl_employees WHERE id = $1", transaction.EmployeeID).Scan(
				&employee.ID,
				&employee.Name,
				&employee.PhoneNumber,
				&employee.Address,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			customer := Customer{}
			err = db.QueryRow("SELECT * FROM tbl_customer WHERE id = $1", transaction.CustomerID).Scan(
				&customer.Id,
				&customer.Name,
				&customer.PhoneNumber,
				&customer.Address,
			)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			rowsBill, err := db.Query("SELECT * FROM bill_details WHERE transaction_id = $1", transaction.Id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer rowsBill.Close()

			billDetails := []BillDetail{}
			for rowsBill.Next() {
				billDetail := BillDetail{}
				err := rowsBill.Scan(
					&billDetail.ID,
					&billDetail.TransactionID,
					&billDetail.ProductID,
					&billDetail.Qty,
				)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				product := Products{}
				err = db.QueryRow("SELECT * FROM tbl_products WHERE id = $1", billDetail.ProductID).Scan(
					&product.ID,
					&product.Name,
					&product.Price,
					&product.Unit,
				)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				billDetail.Product = product
				billDetails = append(billDetails, billDetail)
			}
			transaction.Employee = employee
			transaction.Customer = customer
			transaction.BillDetails = billDetails

			transactions = append(transactions, transaction)
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Transaction detail",
			"data":    transactions,
		})
	}
}
