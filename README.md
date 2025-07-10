<p align="center">
  <img src="./hospos-frontend/public/logo-hospos.svg" alt="HOSPOS Logo" width="180"/>
</p>

# HOSPOS — Modern OPEN SOURCE Restaurant POS System

HOSPOS is a modern, full-stack Point of Sale system for restaurants, cafes, and hospitality venues. It features a robust Go backend and a sleek Next.js/React/Tailwind frontend, designed for speed, reliability, and a beautiful user experience.

---

## Features

- **User/Admin/Role Management**: Secure PIN-based login, role-based access, and robust user controls.
- **Bookings**: Create, view, and manage table bookings with customer assignment, notes, products, bill totals, and status (Open/Closed/Cancelled). Bookings are grouped by date (previous, today, future).
- **Customers**: Add, search, and view customers with full MongoDB persistence and error handling.
- **Products, Categories, Tills/Locations**: Manage your menu, categories, and tills/locations from the admin UI.
- **Discounts**: Static, code, and timed discounts with renew/delete actions and expiry logic.
- **Receipts**: (Planned) Store and view receipts for each booking.
- **CORS & Security**: Secure CORS, error handling, and logging middleware.
- **Modern UI**: Responsive, dark-mode-friendly, and easy to use on touch screens.

---

## Tech Stack

- **Backend**: Go, MongoDB, REST API
- **Frontend**: Next.js, React, Tailwind CSS
- **Database**: MongoDB (collections auto-initialized)

---

## Getting Started

### Prerequisites
- Go 1.20+
- Node.js 18+
- MongoDB (local or cloud)

### Backend
```bash
cd Backend
# Set up your MongoDB connection in internal/db/mongo.go
# Run the server
GO111MODULE=on go run main.go
```

### Frontend
```bash
cd hospos-frontend
npm install
npm run dev
```

- The frontend will be available at [http://localhost:3000](http://localhost:3000)
- The backend API runs on [http://localhost:8080](http://localhost:8080)

### API Proxy
The frontend is configured to proxy `/api/*` requests to the Go backend for seamless integration.

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
hospos-frontend/
  src/app/
    bookings/
    customers/
    ...
  public/
    logo-hospos.svg
```

---

## Customization & Extending
- Add new features (e.g., table map, receipt printing, payment integration) as needed.
- UI and API are modular and easy to extend.

---

## License
MIT

---

<p align="center">
  <img src="./hospos-frontend/public/logo-hospos.svg" alt="HOSPOS Logo" width="120"/>
</p>
