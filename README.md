# **Description**

This project is a Redis-based caching layer for database systems, written in Go. It provides a high-performance, scalable solution to reduce database load by caching frequently accessed data. The system ensures efficient cache management and synchronization with the underlying database.

# **Features**

Redis Integration: Utilizes Redis for fast in-memory data storage.

Automatic Cache Management: Handles cache expiration and invalidation.

High Performance: Reduces database queries by serving data directly from the cache.

Scalability: Supports distributed setups with Redis clusters.

Easy Configuration: Allows customization of cache policies and TTL (Time-To-Live).

# **Requirements**

To build and run this project, you need:

Go: Version 1.18 or newer.

Redis: Version 6.0 or newer.

Database: Any SQL/NoSQL database supported by your application.

# **Installation**

Clone the Repository

```
git clone https://github.com/LOOK-MOM-I-CAN-FLY/DB-cache-Redis.git
cd DB-cache-Redis
```

Set Up Environment Variables

Create a .env file in the project root and configure the following:
```
REDIS_HOST=localhost
REDIS_PORT=6379
DB_CONNECTION_STRING=your_database_connection_string
CACHE_TTL=60 # Cache time-to-live in seconds
```

# **Install** **Dependencies**

Use go mod to install required libraries:

```go mod tidy```

Run the Application

Start the application with the following command:

```go run main.go```

# **Usage**

Key Caching: Cache specific database query results with customizable keys.

Automatic Expiration: Cached data expires based on the configured TTL.

Manual Invalidation: Invalidate specific cache entries when data changes in the database.

# **File Structure**

main.go: Entry point of the application.

cache/: Contains caching logic and Redis integration.

database/: Handles database connections and queries.

config/: Manages configuration and environment variables.

utils/: Utility functions and helpers.

# **Contributing**

Contributions are welcome! If you want to contribute:

1. Fork the repository.

2. Create a new branch for your feature or bug fix:

```git checkout -b feature-name```

3. Commit your changes:

```git commit -m "Description of changes"```

4. Push your branch:

```git push origin feature-name```

5. Create a Pull Request.

# **Future Plans**

Add support for multiple cache eviction policies (e.g., LRU, LFU).

Enhance monitoring and metrics for cache performance.

Implement a write-through cache strategy.

Support advanced Redis features like streams and pub/sub.

Provide examples for integration with popular databases.

# **Acknowledgments**

Redis: For providing a robust in-memory data store.

Open Source Community: For tools and libraries that make this project possible.

Database Systems: For inspiring efficient data access solutions.

# **Contact**

If you have any questions or suggestions, feel free to reach out:

GitHub Issues: Submit your issues here.

Email: igrik315.nekrasov@yandex.ru

Telegram: @Sindi_hall


------------------------------------------------------------------------------

Enjoy using this caching layer and feel free to contribute to its development!
