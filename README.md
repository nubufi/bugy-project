# Buggy Project Fixes

## Identified Issues and Fixes

### 1. Unclosed of Database Connection
- **Problem**: The database connection was opened but never closed, which could lead to connection leaks and resource exhaustion.
- **Fix**: Added `defer db.Close()` in the `main()` function to ensure that the database connection is closed when the application shuts down. Also, added `db.Ping()` to verify that the connection is alive before proceeding.

### 2. Incorrect Use of Goroutines
- **Problem**: Goroutines were unnecessarily used in `getUsers` and `createUser`, and they were immediately blocked by `wg.Wait()`, negating any concurrency benefits. Additionally, there was unsafe concurrent access to `http.ResponseWriter`.
- **Fix**: Removed the unnecessary goroutines. Since each HTTP request runs in its own goroutine, the database operations are now handled synchronously, simplifying the code and avoiding race conditions.

### 3. SQL Injection Vulnerability
- **Problem**: The `createUser` function used string concatenation to build SQL queries, which made the application volnurable to SQL injection attacks. For example an username input as ```John'); DROP TABLE users; --``` could cause deleting the users table.
- **Fix**: Replaced the concatenated SQL query with a parameterized query to prevent SQL injection. Now, the query uses `db.Exec("INSERT INTO users (name) VALUES ($1)", username)` to safely insert the username into the database.

### 4. Inadequate Error Handling
- **Problem**: The code ignored errors from database operations such as `db.Query` and `rows.Scan`, leading to potential silent failures in production.
- **Fix**: Implemented proper error handling for all database operations. If an error occurs, the handler responds with an appropriate HTTP error code and logs the error.
