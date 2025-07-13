import 'package:flutter/material.dart';

class PinLoginScreen extends StatefulWidget {
  final VoidCallback? onLogin;
  const PinLoginScreen({Key? key, this.onLogin}) : super(key: key);

  @override
  State<PinLoginScreen> createState() => _PinLoginScreenState();
}

class _PinLoginScreenState extends State<PinLoginScreen> {
  final _pinController = TextEditingController();
  String? _error;
  bool _isLoading = false;

  Future<void> _login() async {
    setState(() { _error = null; _isLoading = true; });
    if (_pinController.text.length < 3) {
      setState(() { _error = 'PIN must be at least 3 digits.'; _isLoading = false; });
      return;
    }
    await Future.delayed(const Duration(seconds: 1));
    setState(() { _isLoading = false; });
    if (widget.onLogin != null) widget.onLogin!();
    // TODO: Navigate to sales/products screen on success
  }

  @override
  Widget build(BuildContext context) {
    final accentColor = Colors.orangeAccent;
    return Scaffold(
      backgroundColor: Colors.indigo[50],
      body: Center(
        child: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.all(32.0),
            child: Container(
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(24),
                boxShadow: [
                  BoxShadow(
                    color: Colors.indigo.withOpacity(0.15),
                    blurRadius: 16,
                    offset: const Offset(0, 8),
                  ),
                ],
              ),
              child: Padding(
                padding: const EdgeInsets.all(24.0),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Image.asset('assets/logo-hospos.png', height: 80),
                    const SizedBox(height: 24),
                    Text('PIN Login', style: TextStyle(
                      fontSize: 28,
                      fontWeight: FontWeight.bold,
                      color: Colors.indigo[700],
                    )),
                    const SizedBox(height: 24),
                    TextField(
                      controller: _pinController,
                      keyboardType: TextInputType.number,
                      obscureText: true,
                      decoration: InputDecoration(
                        labelText: 'Enter PIN',
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(12),
                        ),
                        prefixIcon: Icon(Icons.lock, color: accentColor),
                      ),
                    ),
                    const SizedBox(height: 20),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton.icon(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.indigo,
                          foregroundColor: Colors.white,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                          ),
                          elevation: 4,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                        ),
                        icon: Icon(Icons.login, color: accentColor),
                        label: _isLoading
                          ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                          : const Text('Login', style: TextStyle(fontSize: 18)),
                        onPressed: _isLoading ? null : _login,
                      ),
                    ),
                    if (_error != null) ...[
                      const SizedBox(height: 12),
                      Text(_error!, style: const TextStyle(color: Colors.red)),
                    ],
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
