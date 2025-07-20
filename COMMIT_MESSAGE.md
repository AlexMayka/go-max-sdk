Initial implementation of Max Bot SDK for Go

- Created comprehensive type system for Max Bot API
  - Models for core entities (Bot, User, Chat, Message, etc.)
  - Request/Response structures for all API endpoints
  - Support for inline keyboards and buttons
  
- Implemented flexible HTTP client with reflection
  - Automatic parsing of struct tags (path, query, json)
  - Context support for request cancellation
  - Built-in authorization via access_token
  - Type-safe endpoint configuration
  
- Established project architecture
  - Clean separation of concerns (types/, internal/)
  - Extensible design for future features (Transport, Router, FSM)
  - Configuration-driven API endpoints

This SDK provides a foundation for building Max messenger bots in Go with a declarative API similar to popular frameworks like Aiogram and Gin.