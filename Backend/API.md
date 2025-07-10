# HOSPOS Backend API Reference

This document describes the main REST API endpoints for the HOSPOS backend.

---

## Authentication
- Most endpoints require a valid user session (PIN-based login).
- Role-based access enforced for admin/user actions.

---

## Endpoints

### Customers
- `GET /api/customers?q=search` — List/search customers by name/email/phone
- `GET /api/customers/{id}` — Get a single customer by ID
- `POST /api/customers` — Add a new customer
- `PATCH /api/customers?id={id}` — Update a customer
- `DELETE /api/customers?id={id}` — Delete a customer

#### Customer Object
```json
{
  "id": "...",
  "name": "...",
  "email": "...",
  "phone": "...",
  "notes": "...",
  "tags": ["..."],
  "createdAt": "..."
}
```

---

### Bookings
- `GET /api/bookings` — List all bookings
- `GET /api/bookings/{id}` — Get booking details
- `POST /api/bookings` — Create a booking
- `PATCH /api/bookings/{id}` — Update status or booking time

#### Booking Object
```json
{
  "id": "...",
  "customerId": "...",
  "tableNumber": "...",
  "products": [
    { "productId": "...", "name": "...", "qty": 1, "price": 9.99 }
  ],
  "billTotal": 0,
  "status": "open|closed|cancelled",
  "createdAt": "...",
  "bookingTime": "...",
  "closedAt": "...",
  "notes": "..."
}
```

---

### Products
- `GET /api/products` — List products
- `GET /api/products/{id}` — Get product details
- `POST /api/products` — Add product
- `PATCH /api/products/{id}` — Update product
- `DELETE /api/products/{id}` — Delete product

---

### Discounts
- `GET /api/discounts` — List discounts
- `POST /api/discounts` — Add discount
- `PATCH /api/discounts/{id}` — Update/renew discount
- `DELETE /api/discounts/{id}` — Delete discount

---

### Users & Auth
- `POST /api/login` — Login with PIN
- `GET /api/users` — List users
- `POST /api/users` — Add user
- `PATCH /api/users/{id}` — Update user
- `DELETE /api/users/{id}` — Delete user

---

## Status Codes
- `200 OK` — Success
- `201 Created` — Resource created
- `400 Bad Request` — Invalid input
- `401 Unauthorized` — Not logged in/invalid PIN
- `403 Forbidden` — Insufficient role
- `404 Not Found` — Resource not found
- `500 Internal Server Error` — Server/database error

---

## Notes
- All endpoints return JSON.
- Timestamps are ISO8601 strings.
- For more details, see the main README or code comments.
