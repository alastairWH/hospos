```markdown
<p align="center">
  <img src="/logo-hospos.svg" alt="HOSPOS Logo" width="180"/>
</p>

# HOSPOS Frontend

This is the modern Next.js/React/Tailwind frontend for the HOSPOS restaurant POS system.

---

## Features
- Responsive, touch-friendly UI for tills, bookings, and admin
- Bookings: create, view, and manage bookings (with customer, table, products, notes, status)
- Customers: add, search, and view customers
- Products, categories, tills/locations management
- Discounts and role-based access
- Modern navigation with logo and dark mode
- API proxy for seamless backend integration

---

## Getting Started

### Prerequisites
- Node.js 18+
- Backend API running (see project root README)

### Setup
```bash
npm install
npm run dev
```
- Open [http://localhost:3000](http://localhost:3000) to use the app.

---

## Project Structure
```
src/app/
  bookings/
  customers/
  ...
public/
  logo-hospos.svg
```

---

## Customization
- Edit `src/app/Navbar.tsx` for navigation and branding
- Add new pages or modals in `src/app/`
- Tailwind CSS for easy UI changes

---

<p align="center">
  <img src="/logo-hospos.svg" alt="HOSPOS Logo" width="120"/>
</p>
```
