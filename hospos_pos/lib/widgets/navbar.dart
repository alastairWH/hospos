import 'package:flutter/material.dart';

class NavBar extends StatelessWidget {
  final String selected;
  final String userRole;
  final VoidCallback? onSales;
  final VoidCallback? onBookings;
  final VoidCallback? onAdmin;
  final VoidCallback? onSupport;
  final VoidCallback? onLogout;

  const NavBar({
    Key? key,
    required this.selected,
    required this.userRole,
    this.onSales,
    this.onBookings,
    this.onAdmin,
    this.onSupport,
    this.onLogout,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        IconButton(
          icon: const Icon(Icons.settings, color: Colors.white),
          onPressed: onSupport,
          tooltip: 'Support',
        ),
        const SizedBox(width: 8),
        _NavBarButton(
          label: 'Sales',
          selected: selected == 'Sales',
          onTap: onSales ?? () {},
        ),
        const SizedBox(width: 8),
        _NavBarButton(
          label: 'Bookings',
          selected: selected == 'Bookings',
          onTap: onBookings ?? () {},
        ),
        const SizedBox(width: 8),
        if (userRole.toLowerCase() == 'admin' || userRole.toLowerCase() == 'manager')
          _NavBarButton(
            label: 'Admin',
            selected: selected == 'Admin',
            onTap: onAdmin ?? () {},
          ),
        const Spacer(),
        IconButton(
          icon: const Icon(Icons.logout, color: Colors.white),
          onPressed: onLogout,
          tooltip: 'Logout',
        ),
      ],
    );
  }
}

class _NavBarButton extends StatelessWidget {
  final String label;
  final bool selected;
  final VoidCallback onTap;
  const _NavBarButton({required this.label, required this.selected, required this.onTap, Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        decoration: BoxDecoration(
          color: selected ? Colors.white : Colors.indigo[700],
          borderRadius: BorderRadius.circular(12),
        ),
        child: Text(
          label,
          style: TextStyle(
            color: selected ? Colors.indigo : Colors.white,
            fontWeight: FontWeight.bold,
            fontSize: 16,
          ),
        ),
      ),
    );
  }
}
