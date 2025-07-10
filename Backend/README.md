<p align="center">
  <img src="../hospos-frontend/public/logo-hospos.svg" alt="HOSPOS Logo" width="180"/>
</p>

# HOSPOS Backend

This is the backend for the HOSPOS Point of Sale system, written in Go. It provides secure, robust RESTful APIs for all core POS features.

---

## Features
- Product and category management
- Table bookings (with customer, table, products, notes, status)
- Customer management (search, add, view)
- Sales with UK VAT integration
- Inventory management
- User, admin, and role management (PIN-based auth, role-based access)
- Reporting and analytics
- Discounts and promotions (static, code, timed)
- Multi-location/till support
- Receipt handling (planned)
- Reservation reminders
- CORS, logging, and error handling middleware
- MongoDB persistence (collections auto-initialized)

---

## Getting Started

### Prerequisites
- Go 1.20+
- MongoDB (local or cloud)

### Setup
```sh
cd Backend
go mod tidy
go run main.go
```
- The backend API runs on [http://localhost:8080](http://localhost:8080)

---

## Project Structure
```
Backend/
  main.go
  internal/
    customers/
    bookings/
    products/
    ...
```

---

## Extending
- Add new feature modules in `internal/`
- All endpoints are modular and easy to extend

---

<p align="center">
  <img src="../hospos-frontend/public/logo-hospos.svg" alt="HOSPOS Logo" width="120"/>
</p>
