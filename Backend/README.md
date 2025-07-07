# HOSPOS Backend

This is the backend for the HOSPOS Point of Sale system, written in Go. It provides RESTful APIs for:

- Product management
- Sales with UK VAT integration
- Table bookings
- Inventory management
- User and role management
- Reporting and analytics
- Customer management
- Payment integration (placeholder)
- Receipt handling
- Discounts and promotions
- Multi-location support
- Offline mode (sync logic placeholder)
- Reservation reminders

## Getting Started

1. Ensure you have Go installed (1.20+ recommended).
2. Clone this repository and run:

```sh
go mod tidy
go run main.go
```

## Project Structure

- `main.go`: Entry point
- `internal/`: Feature modules

## Next Steps
- Implement feature modules in `internal/`
- Add database integration
- Build RESTful endpoints

---

For more details, see `.github/copilot-instructions.md`.
