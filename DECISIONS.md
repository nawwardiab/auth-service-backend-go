# Architecture & Design Decisions

## 1. Why Go + Echo Framework?

**Decision:** Use Go with Echo framework for the backend

**Reasoning:**
- Learned Go during internship - familiar and productive
- Echo provides lightweight routing with powerful middleware
- Built-in validation and error handling
- Performance is excellent for HTTP services
- Similar patterns to Express (easy to explain to TypeScript teams)

**Trade-offs:**
- Go has less ecosystem than Node.js
- Echo is less popular than Gin or standard library
- BUT: Echo's middleware chaining made JWT + CSRF implementation cleaner

**What I'd do differently:**
- For a TypeScript shop, I'd use NestJS or Express with TypeScript
- But the architectural patterns (handler-service-repo) transfer directly

---

## 2. Why Handler-Service-Repository Pattern?

**Decision:** Separate concerns into three layers

**Reasoning:**
- **Handlers**: HTTP concerns only (request/response, validation)
- **Services**: Business logic (can be reused, testable without HTTP)
- **Repositories**: Data access (can swap databases without changing business logic)

**Benefits:**
- Easy to test (can mock each layer)
- Easy to change (swap Echo for Gin, swap Postgres for MySQL)
- Easy to understand (clear responsibilities)

**Inspiration:**
- This is similar to MVC but with clearer separation
- Common in enterprise Go projects
- Makes testing easier (though I haven't added tests yet)

---

## 3. Why JWT in Cookies + CSRF Tokens?

**Decision:** Use JWT stored in HTTPOnly cookies with CSRF protection

**Reasoning:**
- **JWT**: Stateless authentication, no database lookup per request
- **HTTPOnly cookies**: Protected from XSS attacks (JavaScript can't read)
- **CSRF tokens**: Protected from CSRF attacks (cross-site requests blocked)

**Why not sessions?**
- Sessions require database/Redis lookup per request (slower)
- Sessions harder to scale horizontally
- JWT works better for microservices (token is self-contained)

**Why not JWT in localStorage?**
- Vulnerable to XSS attacks
- Can't set HTTPOnly flag
- CSRF protection is easier than XSS protection

**Trade-offs:**
- Can't revoke JWTs easily (would need refresh tokens + blacklist)
- Token size is larger than session ID
- BUT: For this project, security > convenience

**What I'd add in production:**
- Refresh token rotation
- Token blacklist for logout
- Rate limiting

---

## 4. Why TypeScript on Frontend Only?

**Decision:** Backend in Go, frontend in TypeScript

**Reasoning:**
- Built backend during internship when learning Go
- Wanted to demonstrate TypeScript knowledge for the job
- Full-stack TypeScript would be ideal, but I prioritized:
  - Working system over perfect stack match
  - Quality backend over learning Node.js + TypeScript backend rushed

**For the interview:**
- I understand TypeScript/Node.js backend patterns
- NestJS is similar to my Go structure (dependency injection, decorators, layers)
- Concepts transfer (handlers → controllers, repos → services with TypeORM)

---

## 5. Why Address Management as First Feature?

**Decision:** Addresses as the first entity after auth

**Reasoning:**
- **Real-world relevance**: E-commerce sites need shipping addresses
- **Demonstrates relationships**: User has-many Addresses (foreign key)
- **Shows CRUD operations**: Create, Read, Update, Delete
- **Authorization practice**: Users can only see their own addresses

**What this teaches:**
- Data modeling (normalized schema)
- Authorization (ownership checks)
- Full-stack flow (backend CRUD + frontend forms)

---

## 6. Why PostgreSQL Over MySQL/SQLite?

**Decision:** PostgreSQL as the database

**Reasoning:**
- Used PostgreSQL during internship
- Better JSON support than MySQL (useful for future features)
- More robust than SQLite for multi-user apps
- Industry standard for modern web apps
- Excellent migration tools

**Trade-offs:**
- Heavier than SQLite (but fine for production)
- Slightly more complex than MySQL
- BUT: Best choice for scalability

---

## 7. Why Migrations Instead of ORM?

**Decision:** Raw SQL migrations vs. ORM (GORM, etc.)

**Reasoning:**
- **Full control**: Know exactly what SQL runs
- **Transparency**: No hidden queries or N+1 problems
- **Learning**: Forces me to understand database design
- **Performance**: No ORM overhead

**Trade-offs:**
- More verbose (write more SQL)
- No type safety on queries
- BUT: Better for learning and debugging

**Future consideration:**
- For rapid development, I might use an ORM
- For performance-critical code, raw SQL is better
- Hybrid approach: ORM for simple queries, raw SQL for complex

---

## 8. Why No Tests Yet?

**Decision:** Prioritize working features over tests initially

**Reasoning:**
- **Honest answer**: I'm still learning testing best practices
- **Time constraint**: 2 months to build during internship
- **Priority**: Get auth working securely first

**What I know about testing:**
- Unit tests: Test individual functions (handlers, services)
- Integration tests: Test full flow with test database
- E2E tests: Test through frontend

**What I'd add first:**
- Handler tests with mocked services
- Service tests with mocked repos
- At least 70% coverage before production

**Learning plan:**
- Study Go testing package
- Learn testify library for assertions
- Practice TDD on next project

---

## 9. Why This Project Structure?

**Decision:** `internal/` for all application code

**Reasoning:**
- Go convention: `internal/` makes code private to this module
- Clear separation: `internal/handler`, `internal/service`, etc.
- Easy navigation: Everything is organized by responsibility

**Alternative considered:**
- Flat structure (all files in root) → Too messy
- Feature folders (auth/, address/) → Less clear for small projects

---

## 10. Why Custom Validator Wrapper?

**Decision:** Wrap go-playground/validator in custom validator

**Reasoning:**
- Echo requires specific interface
- Centralized validation logic
- Can add custom rules later
- Single place to configure validation

**Shows:**
- Understanding of interfaces
- DRY principle
- Extensibility thinking

---

## Future Decisions (If Asked)

### "How would you add products?"

**Decision:** Products as a separate service

**Structure:**
```
internal/
  product/
    handler.go
    service.go
    repo.go
  model/
    product.go
```

**Why separate:**
- Products don't need auth logic (public catalog)
- Can scale independently
- Clear boundaries

**Integration with auth:**
- Protected routes for admin (add/edit products)
- Public routes for viewing
- Auth service validates tokens, product service handles business logic

### "How would you add cart?"

**Decision:** Cart as session-based (short term) or database-backed (long term)

**Session-based (MVP):**
- Store cart in cookie/session
- Fast, no database
- Lost on logout (acceptable for MVP)

**Database-backed (production):**
- Cart table: user_id, product_id, quantity
- Persists across sessions
- Can track abandoned carts

**My approach:**
- Start with session for speed
- Migrate to database when needed
- Shows I understand MVP vs. production trade-offs

### "How would you add tests?"

**Priority order:**
1. **Service layer tests** (business logic, no HTTP/DB dependencies)
2. **Handler tests** (HTTP layer, mock services)
3. **Integration tests** (full stack with test DB)
4. **Frontend tests** (component + E2E)

**Example test structure:**
```go
// service_test.go
func TestCreateAddress(t *testing.T) {
    mockRepo := &MockAddressRepo{}
    service := NewAddressService(mockRepo)
    
    // Test valid address
    // Test invalid data
    // Test authorization
}
```

---

## Questions I'm Ready For

✅ Why Go?
✅ Why this architecture?
✅ Why JWT + CSRF?
✅ How does auth work?
✅ How would you scale this?
✅ Why no tests?
✅ How would you add feature X?
✅ What would you improve?
✅ Why should we hire you?
```
