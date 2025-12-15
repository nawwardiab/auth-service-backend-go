# Architectural Decisions

This document explains the key technical decisions made in this project and the reasoning behind them.

---

## 1. Why Go + Echo Framework?

**Decision:** Backend built in Go using Echo web framework.

**Reasoning:**
- [Fill in: Why did you choose Go? Was it the internship? Performance? Learning opportunity?]
- [Fill in: Why Echo specifically? Minimalist? Good middleware support?]

**Trade-offs:**
- Pros: [e.g., Fast, compiled, good concurrency, strong typing]
- Cons: [e.g., Company uses TypeScript, less ecosystem than Node.js]

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
Service  → business logic, orchestration, domain rules
Repository → database queries, data access only
```

**Trade-offs:**
- Pros: [e.g., Testable layers, clear responsibilities, easy to swap implementations]
- Cons: [e.g., More files, can feel like boilerplate for simple CRUD]

**Why this matters:** [e.g., Makes adding new features predictable - just follow the pattern]

---

## 3. Why JWT + CSRF Protection Together?

**Decision:** JWT tokens in HTTPOnly cookies + CSRF tokens in separate cookies.

**Reasoning:**
- [Fill in: Why HTTPOnly cookies for JWT? What attack does this prevent?]
- [Fill in: Why CSRF tokens needed? What's the attack vector?]
- [Fill in: Why two separate mechanisms?]

**Implementation:**
- JWT: stored in HTTPOnly cookie (JavaScript can't read it)
- CSRF: stored in readable cookie, sent in request headers
- Backend validates both on protected routes

**Trade-offs:**
- Pros: [e.g., Protects against XSS + CSRF, stateless auth]
- Cons: [e.g., Can't easily revoke JWTs, CSRF cookie readable by JS]

**What I learned:** [Fill in: The CSRF token accessibility challenge you faced - how you solved it]

---

## 4. Why Standardized Error Response System? (Recent Refactor)

**Decision:** All API errors return `{ "error": { "code": string, "message": string } }` format.

**Reasoning:**
- [Fill in: What was wrong with the original error handling?]
- [Fill in: Why does consistent error format matter for API consumers?]
- [Fill in: Why use error codes instead of just messages?]

**Implementation:**
- Created `internal/response` package with error envelope
- Defined 16 typed error code constants (e.g., `AUTH_INVALID_CREDENTIALS`)
- Service layer uses exported error sentinels (`ErrUserNotFound`)
- Handler layer maps service errors to HTTP responses
- Proper dependency flow: handler → service → repo

**Trade-offs:**
- Pros: [e.g., Client-friendly, stable contracts, maintainable, debuggable]
- Cons: [e.g., More upfront work, requires discipline across team]

**Why this shows maturity:** [Fill in: How this demonstrates production thinking, API design awareness]

---

## 5. Why TypeScript on Frontend Only?

**Decision:** Backend in Go, frontend in React + TypeScript.

**Reasoning:**
- [Fill in: Why TypeScript for frontend? Type safety? Better DX?]
- [Fill in: Why not rewrite backend in TypeScript?]

**What this demonstrates:**
- I know TypeScript (frontend proves it)
- I understand the same architectural patterns apply across languages
- Handler-Service-Repo in Go = Controller-Service-Repo in NestJS

**If joining a TypeScript backend team:**
- [Fill in: How would you ramp up? What transfers? What's new?]

---

## 6. Why PostgreSQL?

**Decision:** PostgreSQL for persistence, with migrations via golang-migrate.

**Reasoning:**
- [Fill in: Why relational DB? Why not NoSQL?]
- [Fill in: Why PostgreSQL specifically?]
- [Fill in: Why migrations over ORM?]

**Trade-offs:**
- Pros: [e.g., ACID, relations, strong typing, mature ecosystem]
- Cons: [e.g., Schema changes require migrations, less flexible than NoSQL]

**Migration approach:** [Fill in: Why manual migrations? Control? Learning?]

---

## 7. Why No Tests Yet?

**Decision:** Currently no automated tests (unit, integration, or E2E).

**Honest reasoning:**
- [Fill in: What did you prioritize instead? Security? Architecture?]
- [Fill in: Why was this the right trade-off at the time?]
- [Fill in: What's your actual testing experience level?]

**What I know about testing:**
- I understand the theory: unit tests for services, integration tests for handlers, E2E for flows
- I know what I'd test first: [e.g., service layer business logic, then handler input validation]
- I've seen test structures in the internship codebase

**What I'd do next:**
1. [e.g., Start with service layer tests - pure business logic, easy to mock]
2. [e.g., Add handler tests with mock services]
3. [e.g., Integration tests for critical flows like auth]

**Why I'm eager to learn:** [Fill in: How your team does testing, their practices, getting hands-on experience]

---

## 8. How Would I Add Products/Cart Features?

**Decision:** Deliberately did NOT add e-commerce features before this interview.

**Reasoning:**
- [Fill in: Why would rushing these be risky?]
- [Fill in: What's the trade-off between scope and quality?]

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
- [Fill in: Consistency with existing architecture, predictable patterns]
- [Fill in: Shows systematic thinking vs ad-hoc code]

---

## 9. What Would I Do With More Time?

**Priority 1: Comprehensive Testing**
- [Fill in: What would you test first and why?]

**Priority 2: Better Error Handling**
- [Fill in: What's still rough? Custom error types? Better logging?]

**Priority 3: Frontend Improvements**
- [Fill in: AuthContext? Better loading states? Address UX?]

**Priority 4: Production Readiness**
- [Fill in: Logging? Monitoring? Rate limiting? CI/CD?]

---

## 10. What I'm Most Proud Of

**1. Security Implementation**
- [Fill in: What makes the JWT+CSRF approach sophisticated?]

**2. Error Response Refactor**
- [Fill in: What does this show about your thinking?]

**3. Clean Architecture**
- [Fill in: Why does separation of concerns matter?]

---

## Questions I'm Ready to Answer

- "Walk me through authentication flow" → [Write 3-4 bullet points]
- "Why Go when we use TypeScript?" → [Write 3-4 bullet points]
- "Why no tests?" → [Write your honest answer]
- "What's your biggest challenge?" → [Write about CSRF token accessibility or another real challenge]
- "Why should we hire you?" → [Write 3 key points]

---

## What This Project Demonstrates

[Fill in: What does this project prove about you as a junior developer? Learning speed? Quality focus? Systematic thinking?]
