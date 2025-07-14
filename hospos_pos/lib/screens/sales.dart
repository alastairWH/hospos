import 'package:flutter/material.dart';
import '../services/api_service.dart';
import '../widgets/navbar.dart';

class SalesScreen extends StatefulWidget {
  final String userName;
  final String userRole;
  final VoidCallback? onLogout;
  const SalesScreen({Key? key, this.userName = 'User', this.userRole = 'Role', this.onLogout}) : super(key: key);

  @override
  State<SalesScreen> createState() => _SalesScreenState();
}

class _SalesScreenState extends State<SalesScreen> {
  List<String> _categories = [];
  String? _selectedCategory;
  List<Map<String, dynamic>> _products = [];
  bool _isLoading = true;
  List<Map<String, dynamic>> _cart = [];
  double get _cartSubtotal => _cart.fold(0.0, (sum, item) => sum + (item['price'] * item['qty']));
  double get _cartTax => _cartSubtotal * 0.2 / 1.2; // 20% VAT included
  double get _cartTotal => _cartSubtotal;

  @override
  void initState() {
    super.initState();
    _fetchCategories();
  }

  Future<void> _fetchCategories() async {
    setState(() { _isLoading = true; });
    final categories = await ApiService.getCategories();
    setState(() {
      _categories = categories;
      _isLoading = false;
    });
  }

  Future<void> _fetchProducts(String category) async {
    setState(() { _isLoading = true; _selectedCategory = category; });
    final products = await ApiService.getProducts();
    setState(() {
      _products = products.where((p) => p['category'] == category).toList();
      _isLoading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final accentColor = Colors.orangeAccent;
    return Scaffold(
      backgroundColor: Colors.indigo[50],
      appBar: AppBar(
        backgroundColor: Colors.indigo,
        elevation: 8,
        title: NavBar(
          selected: 'Sales',
          userRole: widget.userRole,
          onSales: () {},      // Add navigation logic here
          onBookings: () {},   // Add navigation logic here
          onAdmin: () {},      // Add navigation logic here
          onSupport: () {},    // Add support logic here
          onLogout: widget.onLogout,
        ),
      ),
      body: Column(
        children: [
          const SizedBox(height: 16),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Row(
                children: [
                  // Categories/Products area (70%)
                  Expanded(
                    flex: 7,
                    child: Container(
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(24),
                        boxShadow: [
                          BoxShadow(
                            color: Colors.indigo.withOpacity(0.12),
                            blurRadius: 16,
                            offset: const Offset(0, 8),
                          ),
                        ],
                      ),
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text('Categories', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.indigo)),
                            const SizedBox(height: 8),
                            SingleChildScrollView(
                              scrollDirection: Axis.horizontal,
                              child: Row(
                                children: _categories.map((cat) => Padding(
                                  padding: const EdgeInsets.symmetric(horizontal: 8.0),
                                  child: ChoiceChip(
                                    label: Text(cat),
                                    selected: _selectedCategory == cat,
                                    onSelected: (selected) {
                                      if (selected) _fetchProducts(cat);
                                    },
                                  ),
                                )).toList(),
                              ),
                            ),
                            const SizedBox(height: 16),
                            if (_selectedCategory != null) ...[
                              Text('Products in $_selectedCategory', style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w600)),
                              const SizedBox(height: 8),
                              Expanded(
                                child: _isLoading
                                    ? const Center(child: CircularProgressIndicator())
                                    : GridView.count(
                                        crossAxisCount: 3,
                                        childAspectRatio: 1.6,
                                        crossAxisSpacing: 16,
                                        mainAxisSpacing: 16,
                                        children: _products.map((product) => ProductCard(
                                          name: product['name'] ?? 'Unnamed',
                                          price: product['price']?.toDouble() ?? 0.0,
                                          category: product['category'] ?? '',
                                          onTap: () {
                                            setState(() {
                                              final idx = _cart.indexWhere((item) => item['id'] == product['id']);
                                              if (idx >= 0) {
                                                _cart[idx]['qty'] += 1;
                                              } else {
                                                _cart.add({
                                                  'id': product['id'],
                                                  'name': product['name'],
                                                  'price': product['price']?.toDouble() ?? 0.0,
                                                  'qty': 1,
                                                });
                                              }
                                            });
                                          },
                                        )).toList(),
                                      ),
                              ),
                            ],
                          ],
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 24),
                  // Cart area (30%)
                  Expanded(
                    flex: 3,
                    child: Container(
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(24),
                        boxShadow: [
                          BoxShadow(
                            color: Colors.indigo.withOpacity(0.12),
                            blurRadius: 16,
                            offset: const Offset(0, 8),
                          ),
                        ],
                      ),
                      child: Padding(
                        padding: const EdgeInsets.all(16.0),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text('Cart', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.indigo)),
                            const SizedBox(height: 12),
                            Expanded(
                              child: _cart.isEmpty
                                  ? const Center(child: Text('Cart is empty'))
                                  : ListView(
                                      children: _cart.map((item) => CartItemCard(
                                        productName: item['name'],
                                        price: item['price'],
                                        qty: item['qty'],
                                        onRemove: () {
                                          setState(() {
                                            final idx = _cart.indexWhere((cartItem) => cartItem['id'] == item['id']);
                                            if (idx >= 0) {
                                              if (_cart[idx]['qty'] > 1) {
                                                _cart[idx]['qty'] -= 1;
                                              } else {
                                                _cart.removeAt(idx);
                                              }
                                            }
                                          });
                                        },
                                        onEditQty: () async {
                                          final idx = _cart.indexWhere((cartItem) => cartItem['id'] == item['id']);
                                          if (idx >= 0) {
                                            int newQty = _cart[idx]['qty'];
                                            final result = await showDialog<int>(
                                              context: context,
                                              builder: (context) {
                                                return AlertDialog(
                                                  title: const Text('Change Quantity'),
                                                  content: Column(
                                                    mainAxisSize: MainAxisSize.min,
                                                    children: [
                                                      Text('Current quantity: ${_cart[idx]['qty']}'),
                                                      TextField(
                                                        keyboardType: TextInputType.number,
                                                        decoration: const InputDecoration(labelText: 'New quantity'),
                                                        onChanged: (val) {
                                                          newQty = int.tryParse(val) ?? newQty;
                                                        },
                                                      ),
                                                    ],
                                                  ),
                                                  actions: [
                                                    TextButton(
                                                      child: const Text('Cancel'),
                                                      onPressed: () => Navigator.of(context).pop(null),
                                                    ),
                                                    ElevatedButton(
                                                      child: const Text('Update'),
                                                      onPressed: () => Navigator.of(context).pop(newQty),
                                                    ),
                                                  ],
                                                );
                                              },
                                            );
                                            if (result != null && result > 0) {
                                              setState(() {
                                                _cart[idx]['qty'] = result;
                                              });
                                            } else if (result == 0) {
                                              setState(() {
                                                _cart.removeAt(idx);
                                              });
                                            }
                                          }
                                        },
                                      )).toList(),
                                    ),
                            ),
                            const SizedBox(height: 12),
                            Text('Subtotal: £${_cartSubtotal.toStringAsFixed(2)}'),
                            Text('VAT (included): £${_cartTax.toStringAsFixed(2)}'),
                            Text('Total: £${_cartTotal.toStringAsFixed(2)}', style: const TextStyle(fontWeight: FontWeight.bold)),
                            const SizedBox(height: 12),
                            ElevatedButton.icon(
                              style: ElevatedButton.styleFrom(
                                backgroundColor: accentColor,
                                foregroundColor: Colors.white,
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                elevation: 4,
                                padding: const EdgeInsets.symmetric(vertical: 12),
                              ),
                              icon: const Icon(Icons.percent),
                              label: const Text('Apply Discount'),
                              onPressed: () {},
                            ),
                            const SizedBox(height: 12),
                            ElevatedButton.icon(
                              style: ElevatedButton.styleFrom(
                                backgroundColor: Colors.indigo,
                                foregroundColor: Colors.white,
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                elevation: 4,
                                padding: const EdgeInsets.symmetric(vertical: 16),
                              ),
                              icon: const Icon(Icons.payment),
                              label: const Text('Take Payment', style: TextStyle(fontSize: 18)),
                              onPressed: () {},
                            ),
                            const SizedBox(height: 12),
                            ElevatedButton.icon(
                              style: ElevatedButton.styleFrom(
                                backgroundColor: Colors.grey[700],
                                foregroundColor: Colors.white,
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                elevation: 4,
                                padding: const EdgeInsets.symmetric(vertical: 16),
                              ),
                              icon: const Icon(Icons.receipt_long),
                              label: const Text('Print Receipt', style: TextStyle(fontSize: 18)),
                              onPressed: () {},
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class ProductCard extends StatelessWidget {
  final String name;
  final double price;
  final String category;
  final VoidCallback? onTap;
  const ProductCard({Key? key, this.name = 'Product', this.price = 9.99, this.category = 'Category', this.onTap}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Material(
      color: Colors.indigo[100],
      borderRadius: BorderRadius.circular(16),
      elevation: 4,
      child: InkWell(
        borderRadius: BorderRadius.circular(16),
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(12.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text(name, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Colors.indigo)),
              const SizedBox(height: 8),
              Text('£${price.toStringAsFixed(2)}', style: const TextStyle(fontSize: 16, color: Colors.indigo)),
            ],
          ),
        ),
      ),
    );
  }
}

class CartItemCard extends StatelessWidget {
  final String productName;
  final double price;
  final int qty;
  final VoidCallback? onRemove;
  final VoidCallback? onEditQty;
  const CartItemCard({Key? key, required this.productName, required this.price, required this.qty, this.onRemove, this.onEditQty}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Card(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 12),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(productName, style: const TextStyle(fontWeight: FontWeight.w500)),
            Text('x$qty'),
            Text('£${price.toStringAsFixed(2)}', style: const TextStyle(fontWeight: FontWeight.bold)),
            IconButton(
              icon: const Icon(Icons.edit, color: Colors.blue),
              tooltip: 'Edit quantity',
              onPressed: onEditQty,
            ),
            IconButton(
              icon: const Icon(Icons.remove_circle, color: Colors.red),
              tooltip: 'Remove',
              onPressed: onRemove,
            ),
          ],
        ),
      ),
    );
  }
}