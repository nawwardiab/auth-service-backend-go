# Architectural Decisions

This document explains the key technical decisions made in this project and the reasoning behind them.

---

## 1. Why Go + Echo Framework?

**Decision:** Backend built in Go using Echo web framework.

**Reasoning:**
Go was the decision of the CTO in the startup I did internship with. As an experienced SRE his concerns were performance, type-safety and for it was fast and produced an executable (compiled).
Echo has an easy and light-weight HTTP-handling approach, great middleware support (for logging, validating, jwt and cors)

**Trade-offs:**

- Pros: Fast, compiled, strong typing
- Cons: less ecosystem than Node.js, (it required higher entrypoint with complex documentation/steeper learning curve)

**What I'd do differently:** [e.g., If starting fresh for this company, I'd use NestJS/TypeScript to match their stack, but the architectural patterns transfer directly]

---

## 2. Why Handler-Service-Repository Pattern?

**Decision:** Three-layer architecture with clear separation between HTTP, business logic, and data access.

**Reasoning:**
- [Fill in: Why separate these layers? Testing? Maintainability? Clear contracts?]
- [Fill in: How does this help when building features?]

**Structure:**

```
Handler  → validates input, calls service, formats response
Service  → business logic, domain rules
Repository → database queries, data access only
```

**Trade-offs:**

- Pros: Testable layers, clear responsibilities, easy to swap implementations
- Cons: More files

**Why this matters:** [e.g., Makes adding new features predictable - just follow the pattern]

---

## 3. Why JWT + CSRF Protection Together?

**Decision:** JWT tokens in HTTPOnly cookies + CSRF tokens in separate cookies.

**Reasoning:**

- **Why HTTPOnly cookies for JWT?** They cannot be accessed by JavaScript, protecting against XSS (Cross-Site Scripting) attacks where malicious JavaScript tries to steal tokens.
- **Why CSRF tokens needed?** To protect against CSRF (Cross-Site Request Forgery) attacks where an attacker tricks a user's browser into making unwanted requests to our API. The CSRF token verifies that requests actually come from our frontend.
- **Why two separate mechanisms?** Defense-in-depth. Even if an attacker tricks the user's browser into making a request (which includes the JWT cookie), they cannot supply the correct CSRF token because they can't read it from a different origin.

**Implementation:**

- JWT: stored in HTTPOnly cookie (JavaScript can't read it)
- CSRF: stored in readable cookie, sent in request headers
- Backend validates both on protected routes

**Trade-offs:**

- Pros: Protects against both XSS and CSRF attacks
- Cons: Can't easily revoke JWTs, CSRF cookie readable by JS

**What I learned:** [Fill in: The CSRF token accessibility challenge you faced - how you solved it]

---

## 4. Why Standardized Error Response System? (Recent Refactor)

**Decision:** All API errors return `{ "error": { "code": string, "message": string } }` format.

**Reasoning:**

- **What was wrong?** It was inconsistent and not global. Each layer had its own error handling, making it difficult to track which error messages were for clients vs internal (dev) use.
- **Why does consistent format matter?** Clients can programmatically handle errors using stable codes. Easier debugging - we know exactly which business rule failed.
- **Why use error codes instead of just messages?** To better identify issues, improves debugging and maintainability. Error messages might change, but codes stay stable.

**Implementation:**

- Created `internal/response` package with error envelope
- Defined 16 typed error code constants (e.g., `AUTH_INVALID_CREDENTIALS`)
- Service layer uses exported error sentinels (`ErrUserNotFound`)
- Handler layer maps service errors to HTTP responses
- Proper dependency flow: handler → service → repo

**Trade-offs:**

- Pros: Client-friendly, maintainable, debuggable
- Cons: More upfront work, requires discipline across team

**Why this shows maturity:**

- I didn't just fix bugs - I identified systemic technical debt
- I thought about API consumers (clients need stable error codes)
- I refactored methodically (proper dependency flow, separation of concerns)
- I documented the change (shows I think about team communication)
- Real production systems need consistent error contracts for monitoring/alerting

---

## 5. Why TypeScript on Frontend Only?

**Decision:** Backend in Go, frontend in React + TypeScript.

**Reasoning:**

- **Why TypeScript for frontend?** Company decision at internship. I used it mainly to test backend. TypeScript has great type safety, especially when combined with Go's strong typing on the backend.
- **Why not rewrite backend in TypeScript?** I wanted to demonstrate my architectural and engineering thinking for this interview. Providing working, clean code to demonstrate my abilities. Rewriting it in TypeScript wouldn't make me feel confident showing my real abilities, since it would be a completely new language for backend.

**If joining a TypeScript backend team:**
I have experience with Express.js and Node.js. I would leverage this knowledge to scaffold a similar architectural pattern, while learning TypeScript-specific features to get the best out of the language.

---

## 7. Why No Tests Yet?

**Decision:** Currently no automated tests (unit, integration, or E2E).

**Honest reasoning:**

- **What did I prioritize?** Learning the language, security implementation, and clean architecture.
- **Why was this the right trade-off?** I was/am still learning Go's syntax and understanding the language's benefits. My next step is to learn writing tests in Go.
- **Actual testing experience level:** Beginner level - I understand theory but lack hands-on practice. This is my #1 learning priority.

**What I know about testing:**

- I understand the theory: unit tests for services, integration tests for handlers, E2E for flows
- I know what I'd test first: [e.g., service layer business logic, then handler input validation]
- I've seen test structures in the internship codebase

**What I'd do next:**

1. Learn how to write tests in Go (table-driven tests, mocking patterns)
2. Start with service layer tests - pure business logic, easy to mock
3. Add handler tests with mock services
4. Integration tests for critical flows like auth

**Why I'm eager to learn:**
Testing is essential for confidence in production. I want to learn testing practices - what one tests, how one structures tests, and how one balances coverage with development speed. I learn best by doing, so I'm excited to get hands-on experience with real production testing patterns.

---

## 8. How Would I Add Products/Cart Features?

**Decision:** Deliberately did NOT add e-commerce features before this interview.

**Reasoning:**

- **Why would rushing be risky?** I have limited testing experience - rushing features would introduce bugs I couldn't catch in time.
- **Time constraint:** Few time before interview - not enough time to build and test properly.
- **Better strategy:** Show polished, working auth system rather than half-broken e-commerce features.

Rushing e-commerce features would have shown quantity but questioned quality.

**If I were to add them:**

### Products Service

- Same handler-service-repo pattern
- `Product` model: id, name, description, price, stock, created_at
- Public endpoints: GET /products, GET /products/:id
- Protected admin endpoints: POST /products, PATCH /products/:id
- PostgreSQL table with migration

### Cart Service

- Associate cart items with user_id from JWT
- `CartItem` model: id, user_id, product_id, quantity
- Protected endpoints: GET/POST/PATCH/DELETE /cart/items
- Business rules in service layer (stock validation, price calculations)

### Why this approach?

- **Consistency:** Same handler-service-repo pattern as auth/address
- **Predictability:** Any dev can understand it because it follows established patterns
- **Systematic thinking:** I don't reinvent the wheel for each feature

---

## 9. What Would I Do With More Time?

**Priority 1: Comprehensive Testing**

I would start with service layer (pure business logic, no http/db mocking needed. Then handler tests (mock services). Then critical flow integration tests (auth, CRUD).
Testing first, because it is my biggest gap.

**Priority 2: Better Error Handling**

Already standardized, but could improve:

- structured logging
- better error messages for validation failures
- Monitoring based on error codes

**Priority 3: Frontend Improvements**

- Auth context to centralize auth state (reduce API calls)
- Better loading/error states (user feedback)
- Address UX improvements (clearer default address selection)

**Priority 4: Production Readiness**

- Structured logging (JSON logs for parsing)
- Health check endpoints
- Rate limiting (prevent abuse)
- CI/CD pipeline (automated deploy)

---

## 10. What I'm Most Proud Of

**1. Security Implementation**
JWT + CSRF work together for defense-in-depth:

- **JWT (in HTTPOnly cookie):** Authenticates the user, protected from XSS because JavaScript cannot access HTTPOnly cookies
- **CSRF token (readable cookie + header):** Validates that requests come from my frontend, not a malicious site. Even if an attacker tricks the browser into sending the JWT cookie, they cannot provide the correct CSRF token.

Two separate mechanisms protect against two different attack vectors.

**2. Error Response Refactor**
Shows thinking like an engineer. Standardizing errors makes it easier to debug and understand where errors are coming from. It shows I think about maintainability and API design - not just making features work, but making them work _well_ for the long term.

**3. Clean Architecture**
- [Fill in: Why does separation of concerns matter?]

---

## Questions I'm Ready to Answer

**"Walk me through authentication flow"**

1. User submits email/password to /api/login
2. Backend validates credentials (bcrypt compare)
3. If valid: generate JWT with user_id, set in HTTPOnly cookie
4. Also set CSRF token in readable cookie
5. Frontend stores CSRF token, sends in X-CSRF-Token header
6. Protected routes validate both: JWT (who you are) + CSRF (request is legitimate)

**"Why no tests?"**
I prioritized learning Go, security implementation, and clean architecture first. Testing is my next learning priority. I understand the theory (unit → integration → E2E) but lack hands-on experience. I'm eager to learn your team's testing practices.

**"What's your biggest challenge?"**
CSRF cookie configuration in development. Initially set Secure=true everywhere, which broke local testing (no HTTPS). Had to learn about environment-aware configuration - checking ENV variable to set Secure=false in dev, true in prod. Taught me that security configs need to be environment-aware.

**"Why should we hire you?"**

1. **Fast learner:** Go + PostgreSQL from scratch in 2 months
2. **Quality-focused:** Could have rushed e-commerce features, chose to polish auth system instead
3. **Honest about gaps:** I know testing is weak, documentation strong. I'm coachable and eager to grow
4. **Creative problem solver:** [Story about handling the CSRF environment config challenge]

---

## What This Project Demonstrates

This project proves I can:

- Learn new stacks quickly (Go + PostgreSQL in 2 months)
- Think systematically (consistent architecture patterns)
- Prioritize quality over quantity (polished auth vs rushed features)
- Identify and fix technical debt (error response refactor)
- Be honest about gaps and eager to fill them (testing, production ops)
