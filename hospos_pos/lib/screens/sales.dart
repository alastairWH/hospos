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
  List<Map<String, dynamic>> _discounts = [];
  bool _discountsLoading = false;
  String? _selectedDiscountId;
  double _selectedDiscountPercent = 0.0;
  String _selectedDiscountName = '';
  List<String> _categories = [];
  String? _selectedCategory;
  List<Map<String, dynamic>> _products = [];
  bool _isLoading = true;
  List<Map<String, dynamic>> _cart = [];

  // Payment breakdown
  double _cashPaid = 0.0;
  double _cardPaid = 0.0;

  double get _cartSubtotal => _cart.fold(0.0, (sum, item) => sum + (item['price'] * item['qty']));
  double get _cartTax => _cartSubtotal * 0.2 / 1.2; // 20% VAT included

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

  Future<void> _fetchDiscounts() async {
    setState(() { _discountsLoading = true; });
    try {
      final discounts = await ApiService.getDiscounts();
      setState(() {
        _discounts = discounts;
        _discountsLoading = false;
      });
    } catch (e) {
      setState(() { _discountsLoading = false; });
    }
  }

  Future<void> _showDiscountModal() async {
    await _fetchDiscounts();
    final selected = await showDialog<Map<String, dynamic>>(
      context: context,
      builder: (context) {
        return AlertDialog(
          title: const Text('Select Discount'),
          content: _discountsLoading
              ? const Center(child: CircularProgressIndicator())
              : SizedBox(
                  width: 300,
                  child: _discounts.isEmpty
                      ? const Text('No discounts available')
                      : ListView.builder(
                          shrinkWrap: true,
                          itemCount: _discounts.length,
                          itemBuilder: (context, idx) {
                            final discount = _discounts[idx];
                            return ListTile(
                              title: Text(discount['name'] ?? 'Discount'),
                              subtitle: Text('${discount['percent'] ?? 0}% off'),
                              trailing: _selectedDiscountId == discount['id']
                                  ? const Icon(Icons.check, color: Colors.green)
                                  : null,
                              onTap: () {
                                Navigator.of(context).pop(discount);
                              },
                            );
                          },
                        ),
                ),
          actions: [
            TextButton(
              child: const Text('Close'),
              onPressed: () => Navigator.of(context).pop(),
            ),
          ],
        );
      },
    );
    if (selected != null) {
      setState(() {
        _selectedDiscountId = selected['id'] ?? selected['_id'];
        _selectedDiscountPercent = (selected['percent'] ?? 0).toDouble();
        _selectedDiscountName = selected['name'] ?? '';
      });
    }
  }

  Future<void> _showPaymentModal() async {
    await _showCashPaymentModal();
  }

  Future<void> _showCashPaymentModal() async {
    double amountDue = _cartSubtotal - (_selectedDiscountId != null ? _cartSubtotal * (_selectedDiscountPercent / 100) : 0);
    String enteredAmount = ''; // Start empty
    double cashPaid = 0.0;
    double cardPaid = 0.0;
    String paymentType = 'cash'; // default

    await showDialog(
      context: context,
      builder: (context) {
        return StatefulBuilder(
          builder: (context, setModalState) {
            double entered = double.tryParse(enteredAmount) ?? 0.0;
            double change = paymentType == 'cash' && entered > amountDue ? entered - amountDue : 0.0;
            bool showRemaining = enteredAmount.isEmpty || entered <= amountDue;
            return Dialog(
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
              child: Container(
                width: 480,
                padding: const EdgeInsets.all(0),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(18),
                  boxShadow: [BoxShadow(color: Colors.indigo.withOpacity(0.08), blurRadius: 18, offset: const Offset(0, 6))],
                ),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    // Header
                    Container(
                      width: double.infinity,
                      padding: const EdgeInsets.symmetric(vertical: 18, horizontal: 24),
                      decoration: BoxDecoration(
                        color: Colors.indigo,
                        borderRadius: const BorderRadius.only(topLeft: Radius.circular(18), topRight: Radius.circular(18)),
                      ),
                      child: const Text('Payment', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.white)),
                    ),
                    Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Row(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // Left: Summary panel
                          Expanded(
                            flex: 4,
                            child: Container(
                              padding: const EdgeInsets.all(12),
                              decoration: BoxDecoration(
                                color: Colors.indigo[50],
                                borderRadius: BorderRadius.circular(12),
                              ),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text('Amount Due', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Colors.indigo)),
                                  Text('£${amountDue.toStringAsFixed(2)}', style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold, color: Colors.indigo)),
                                  const SizedBox(height: 10),
                                  Text('Tax', style: TextStyle(color: Colors.grey[700], fontWeight: FontWeight.w500, fontSize: 13)),
                                  Text('£${_cartTax.toStringAsFixed(2)}', style: const TextStyle(fontSize: 15)),
                                  const SizedBox(height: 6),
                                  Text('Basket Discount', style: TextStyle(color: Colors.grey[700], fontWeight: FontWeight.w500, fontSize: 13)),
                                  Text('£${_selectedDiscountId != null ? (_cartSubtotal * (_selectedDiscountPercent / 100)).toStringAsFixed(2) : '0.00'}', style: const TextStyle(fontSize: 15)),
                                  const SizedBox(height: 10),
                                  Text('Paid by Cash: £${cashPaid.toStringAsFixed(2)}', style: const TextStyle(fontSize: 15, color: Colors.green)),
                                  Text('Paid by Card: £${cardPaid.toStringAsFixed(2)}', style: const TextStyle(fontSize: 15, color: Colors.blue)),
                                  if (paymentType == 'cash' && entered > amountDue)
                                    Text('Change: £${change.toStringAsFixed(2)}', style: const TextStyle(fontSize: 16, color: Colors.green, fontWeight: FontWeight.bold)),
                                  if (showRemaining && amountDue > 0)
                                    Text('Remaining: £${amountDue.toStringAsFixed(2)}', style: const TextStyle(fontSize: 15, color: Colors.orange, fontWeight: FontWeight.w600)),
                                ],
                              ),
                            ),
                          ),
                          const SizedBox(width: 16),
                          // Right: Numpad and quick cash
                          Expanded(
                            flex: 6,
                            child: Column(
                              children: [
                                // Amount entered box
                                Container(
                                  width: double.infinity,
                                  margin: const EdgeInsets.only(bottom: 8),
                                  padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 12),
                                  decoration: BoxDecoration(
                                    color: Colors.indigo[100],
                                    borderRadius: BorderRadius.circular(8),
                                  ),
                                  child: Text(
                                    'Amount Entered: £${enteredAmount.isEmpty ? '0.00' : enteredAmount}',
                                    style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.indigo),
                                  ),
                                ),
                                // Numpad
                                GridView.count(
                                  crossAxisCount: 3,
                                  shrinkWrap: true,
                                  mainAxisSpacing: 8,
                                  crossAxisSpacing: 8,
                                  childAspectRatio: 1.1,
                                  physics: const NeverScrollableScrollPhysics(),
                                  children: [
                                    ...['1','2','3','4','5','6','7','8','9','.','0','⌫'].map((n) => ElevatedButton(
                                      style: ElevatedButton.styleFrom(
                                        backgroundColor: Colors.indigo[100],
                                        foregroundColor: Colors.indigo[900],
                                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
                                        elevation: 2,
                                        padding: const EdgeInsets.symmetric(vertical: 12),
                                      ),
                                      onPressed: () {
                                        setModalState(() {
                                          if (n == '⌫') {
                                            if (enteredAmount.isNotEmpty) {
                                              enteredAmount = enteredAmount.substring(0, enteredAmount.length - 1);
                                            }
                                          } else if (n == '.') {
                                            if (!enteredAmount.contains('.')) enteredAmount += '.';
                                          } else {
                                            enteredAmount += n;
                                          }
                                        });
                                      },
                                      child: Text(n, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                                    )),
                                  ],
                                ),
                                const SizedBox(height: 8),
                                // Quick cash buttons
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                                  children: [5, 10, 20, 0].map((amt) =>
                                    ElevatedButton(
                                      style: ElevatedButton.styleFrom(
                                        backgroundColor: amt == 0 ? Colors.grey[300] : Colors.orangeAccent,
                                        foregroundColor: amt == 0 ? Colors.black : Colors.white,
                                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
                                        elevation: 2,
                                        padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 12),
                                      ),
                                      onPressed: () {
                                        setModalState(() {
                                          if (amt == 0) {
                                            enteredAmount = amountDue.toStringAsFixed(2);
                                          } else {
                                            enteredAmount = amt.toStringAsFixed(2);
                                          }
                                        });
                                      },
                                      child: Text(amt == 0 ? 'Exact' : '£$amt', style: const TextStyle(fontSize: 15, fontWeight: FontWeight.w600)),
                                    )
                                  ).toList(),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 8),
                    // Pay By buttons
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 8),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                        children: [
                          Expanded(
                            child: ElevatedButton.icon(
                              style: ElevatedButton.styleFrom(
                                backgroundColor: paymentType == 'cash' ? Colors.green : Colors.grey[300],
                                foregroundColor: paymentType == 'cash' ? Colors.white : Colors.black,
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                elevation: 4,
                                padding: const EdgeInsets.symmetric(vertical: 12),
                              ),
                              icon: const Icon(Icons.money),
                              label: const Text('Cash', style: TextStyle(fontSize: 15, fontWeight: FontWeight.bold)),
                              onPressed: () async {
                                setModalState(() {
                                  double entered = double.tryParse(enteredAmount) ?? 0.0;
                                  if (entered > 0 && amountDue > 0) {
                                    double pay = entered > amountDue ? amountDue : entered;
                                    cashPaid += pay;
                                    amountDue -= pay;
                                    enteredAmount = amountDue > 0 ? amountDue.toStringAsFixed(2) : '';
                                  }
                                  paymentType = 'cash';
                                });
                                if (amountDue <= 0) {
                                  if (double.tryParse(enteredAmount) != null && double.tryParse(enteredAmount)! > amountDue) {
                                    double change = double.tryParse(enteredAmount)! - amountDue;
                                    await _showChangeModal(change);
                                  }
                                  Navigator.of(context).pop({'cash': cashPaid, 'card': cardPaid});
                                }
                              },
                            ),
                          ),
                          const SizedBox(width: 16),
                          Expanded(
                            child: ElevatedButton.icon(
                              style: ElevatedButton.styleFrom(
                                backgroundColor: paymentType == 'card' ? Colors.blue : Colors.grey[300],
                                foregroundColor: paymentType == 'card' ? Colors.white : Colors.black,
                                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                elevation: 4,
                                padding: const EdgeInsets.symmetric(vertical: 12),
                              ),
                              icon: const Icon(Icons.credit_card),
                              label: const Text('Card', style: TextStyle(fontSize: 15, fontWeight: FontWeight.bold)),
                              onPressed: () async {
                                setModalState(() {
                                  double entered = double.tryParse(enteredAmount) ?? 0.0;
                                  if (entered > 0 && amountDue > 0) {
                                    double pay = entered > amountDue ? amountDue : entered;
                                    cashPaid += pay;
                                    amountDue -= pay;
                                    enteredAmount = amountDue > 0 ? amountDue.toStringAsFixed(2) : '';
                                  }
                                  paymentType = 'cash';
                                });
                                if (amountDue <= 0) {
                                  if (double.tryParse(enteredAmount) != null && double.tryParse(enteredAmount)! > amountDue) {
                                    double change = double.tryParse(enteredAmount)! - amountDue;
                                    await _showChangeModal(change); // This must be awaited BEFORE closing the payment modal
                                  }
                                  Navigator.of(context).pop({'cash': cashPaid, 'card': cardPaid});
                                }
                              },
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(height: 10),
                  ],
                ),
              ),
            );
          },
        );
      },
    ).then((result) {
      if (result is Map) {
        setState(() {
          _cashPaid = result['cash'] ?? 0.0;
          _cardPaid = result['card'] ?? 0.0;
          // Clear cart for new sale
          _cart.clear();
          _selectedDiscountId = null;
          _selectedDiscountPercent = 0.0;
          _selectedDiscountName = '';
        });
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Payment complete: Cash £${_cashPaid.toStringAsFixed(2)}, Card £${_cardPaid.toStringAsFixed(2)}'),
            backgroundColor: Colors.indigo,
          ),
        );
      }
    });
  }

  // Timer modal for change
  Future<void> _showChangeModal(double change) async {
    await showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) {
        return _ChangeTimerDialog(change: change);
      },
    );
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
          onSales: () {},
          onBookings: () {},
          onAdmin: () {},
          onSupport: () {},
          onLogout: widget.onLogout,
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: Padding(
              padding: const EdgeInsets.all(24.0),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
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
                            Text(
                              _selectedDiscountId != null && _selectedDiscountPercent > 0
                                  ? 'Discount: -£${(_cartSubtotal * (_selectedDiscountPercent / 100)).toStringAsFixed(2)} (${_selectedDiscountPercent.toStringAsFixed(0)}%)'
                                  : 'Discount: None applied',
                              style: TextStyle(
                                color: _selectedDiscountId != null && _selectedDiscountPercent > 0 ? Colors.green : Colors.grey[700],
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            Text(
                              'Total: £${(_cartSubtotal - (_selectedDiscountId != null ? _cartSubtotal * (_selectedDiscountPercent / 100) : 0)).toStringAsFixed(2)}',
                              style: const TextStyle(fontWeight: FontWeight.bold),
                            ),
                            const SizedBox(height: 12),
                            Row(
                              children: [
                                Expanded(
                                  flex: _selectedDiscountId != null && _selectedDiscountPercent > 0 ? 5 : 10,
                                  child: ElevatedButton.icon(
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor: accentColor,
                                      foregroundColor: Colors.white,
                                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                      elevation: 4,
                                      padding: const EdgeInsets.symmetric(vertical: 12),
                                    ),
                                    icon: const Icon(Icons.percent),
                                    label: const Text('Apply Discount'),
                                    onPressed: _showDiscountModal,
                                  ),
                                ),
                                if (_selectedDiscountId != null && _selectedDiscountPercent > 0)
                                  const SizedBox(width: 12),
                                if (_selectedDiscountId != null && _selectedDiscountPercent > 0)
                                  Expanded(
                                    flex: 5,
                                    child: ElevatedButton.icon(
                                      style: ElevatedButton.styleFrom(
                                        backgroundColor: Colors.red,
                                        foregroundColor: Colors.white,
                                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                        elevation: 4,
                                        padding: const EdgeInsets.symmetric(vertical: 12),
                                      ),
                                      icon: const Icon(Icons.close),
                                      label: const Text('Remove'),
                                      onPressed: () {
                                        setState(() {
                                          _selectedDiscountId = null;
                                          _selectedDiscountPercent = 0.0;
                                          _selectedDiscountName = '';
                                        });
                                      },
                                    ),
                                  ),
                              ],
                            ),
                            const SizedBox(height: 12),
                            SizedBox(
                              width: double.infinity,
                              child: ElevatedButton.icon(
                                style: ElevatedButton.styleFrom(
                                  backgroundColor: Colors.indigo,
                                  foregroundColor: Colors.white,
                                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                  elevation: 4,
                                  padding: const EdgeInsets.symmetric(vertical: 16),
                                ),
                                icon: const Icon(Icons.payment),
                                label: const Text('Take Payment', style: TextStyle(fontSize: 18)),
                                onPressed: _showPaymentModal,
                              ),
                            ),
                            const SizedBox(height: 12),
                            SizedBox(
                              width: double.infinity,
                              child: ElevatedButton.icon(
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

// Timer modal for change
class _ChangeTimerDialog extends StatefulWidget {
  final double change;
  const _ChangeTimerDialog({Key? key, required this.change}) : super(key: key);

  @override
  State<_ChangeTimerDialog> createState() => _ChangeTimerDialogState();
}

class _ChangeTimerDialogState extends State<_ChangeTimerDialog> {
  double progress = 0.0;

  @override
  void initState() {
    super.initState();
    _startTimer();
  }

  void _startTimer() async {
    while (progress < 1.0) {
      await Future.delayed(const Duration(milliseconds: 30));
      if (!mounted) return;
      setState(() {
        progress += 0.02;
      });
    }
    if (mounted) Navigator.of(context).pop();
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(18)),
      child: Container(
        width: 320,
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            const Text('Change Due', style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold, color: Colors.green)),
            const SizedBox(height: 16),
            Text('£${widget.change.toStringAsFixed(2)}', style: const TextStyle(fontSize: 32, fontWeight: FontWeight.bold, color: Colors.green)),
            const SizedBox(height: 24),
            LinearProgressIndicator(value: progress, minHeight: 8, backgroundColor: Colors.grey[300], color: Colors.green),
            const SizedBox(height: 12),
            const Text('Please give change to customer...', style: TextStyle(fontSize: 16)),
          ],
        ),
      ),
    );
  }
}

// Top-level widget classes
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